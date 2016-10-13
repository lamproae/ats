package cli

import (
	"log"
	"telnet"
//	"expect"
	"errors"
	"product"
//	"os"
	"time"
	"regexp"
)


const timeout = 10 * time.Second

type Cli struct {
	c *telnet.Client
	cMode string	/* Current Mode */
	state bool	/* Current state */
	err   error
	command Command /*Last run command */
	result []byte/* Result of Last command */
	product *product.Product
}

type Expected struct {
	Command Command	`command`
	Expected string	`expected`
	ExpectedNo string `expectedNo`
}

func (e *Expected) GetCommand () string {
	return e.Command.Command
}

func (e *Expected) GetExpected() string {
	return e.Expected
}

/* Maybe Tree like structure is better.*/
var ModeSwitchCommand = map[string][]Command {
	"enable->config" : []Command{{"enable", "config terminal","#"}},
	"config->enable" : []Command{{"config", "exit", "#"}},
	"config->bridge" : []Command{{"config", "bridge", "#"}},
	"bridge->config" : []Command{{"bridge", "exit", "#"}},
	"enable->bridge" : []Command{{"enable", "config terminal", "#"}, {"config", "bridge", "#"}},
	"bridge->enable" : []Command{{"bridge", "exit", "#"}, {"config", "exit", "#"}},
	"interface->enable" : []Command{{"interface ", "exit", "#"}},
	"ospf->enable" : []Command{{"ospf", "exit", "#"}},
	"bgp->enable" : []Command{{"bgp", "exit", "#"}},
}

/* This should check the result of hostname setting and build dynamically. */
var cliModePromptDB map[string]string

type cliMode struct {
	mode string
	prompt string
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln("Error:", err)
	}
}

type Command struct {
	Mode string	`mode`
	Command string	`command`
	End string	`end`
}

func NewCommand (mode, cmd, end string) *Command {
	return &Command{
		Mode : mode,
		Command : cmd,
		End : end,
	}
}

func (c *Cli) Expect(e string) error {
	checkErr(c.c.SetReadDeadline(time.Now().Add(timeout)))
	data, err := c.c.ReadUntil(e)
	checkErr(err)
	//Here we should remove the header and footer of result.
	c.result = data[len(c.command.Command)+2:len(data)-len(c.command.Mode+"#")]
	//log.Println("Command result:")
	//log.Println(string(c.result))
	return nil
}

func (c *Cli) Sendln(l string) error {
	checkErr(c.c.SetWriteDeadline(time.Now().Add(timeout)))
	buf := make([]byte, len(l)+1)
	copy(buf, l)
	buf[len(l)] = '\n'
	_, err := c.c.Write(buf)
	checkErr(err)
	//log.Println("Send command: ", l)
	return nil
}

func (c *Cli) CommandResult() string {
	return string(c.result)
}

func (c *Cli) SwitchMode(mode string) error {
	modeswitch := c.cMode+"->"+mode
	log.Println(modeswitch)
	commands, ok := ModeSwitchCommand[modeswitch] 
	if!ok {
		log.Println("Unkonwn mode switch: ", modeswitch)
		return errors.New("Unkonw mode switch")
	}

	for _, command := range commands {
		c.Sendln(command.Command)
		c.Expect(command.End)
		/* May be we should call RunCommand, but there will be conflict opertion. Check in the future. */
		//c.RunCommand(command)
	}

	c.cMode = mode
	log.Println("CLI mode switch to: ", c.cMode)
	return nil
	/*
	if c.cMode == mode {
		return nil
	} else {
		log.Println("Mode switch failed, current mode is: ", c.cMode)
		return errors.New("Mode switch failed")
	}
	*/
}

func (c *Cli) RunCommand(command Command) {
	tMode := command.Mode
	if c.cMode != tMode {
		c.SwitchMode(tMode)
	}
	log.Println(command.Command)
	log.Println(command.End)
	c.Sendln(command.Command)
	//Test method.
	//data, _ := c.c.ReadBytes('#')
	//os.Stdout.Write(data)
	c.Expect(command.End)
	c.command = command
}

func (c *Cli) Assert (e Expected) (string, error) {
	var result []byte
	c.RunCommand(e.Command)
	if e.Expected != "" {
		expectRe := regexp.MustCompile(e.Expected)
		result := expectRe.Find(c.result)
		if result != nil {
			return "", nil
		}
	}

	if e.ExpectedNo != "" {
		expectNoRe := regexp.MustCompile(e.ExpectedNo)
		result := expectNoRe.Find(c.result)
		if result == nil {
			return "", nil
		}

	}

	return string(result), errors.New("Unexpected")
}

func (c *Cli) Login() error {
	c.Expect(c.product.Hostname()+" login")
	c.Sendln(c.product.Username())
	c.Expect("Password: ")
	c.Sendln(c.product.Password())
	c.Expect(c.product.Hostname()+">")
	c.Sendln("enable")
	c.Expect(c.product.Hostname()+"#")
	return nil
}

func (c *Cli) SetHostname(name string) {
	mode := c.cMode
	c.RunCommand(Command{"enable", "hostname "+name, "#"})
	c.SwitchMode(mode)
}

func New(name string) *Cli {
	products := product.ProductDB()
	p, ok := products[name]
	if !ok {
		log.Println("Unkonw device: ", name)
		return nil
	}

	dst := p.Address() + ":23"
	client, err := telnet.NewClient(dst)
	checkErr(err)
	client.SetUnixWriteMode(true)

	c := Cli {
		c : client,
		state : true,
		err : nil,
		result : nil,
		product : p,
	}

	c.Login()
	c.cMode = p.Initmode()
	c.RunCommand(Command{"enable", "terminal length 0", "#"})

	return &c
}

