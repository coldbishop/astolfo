package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"unicode/utf8"

	"github.com/atotto/clipboard"
	"github.com/spf13/pflag"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	defaultPassLength  = 16
	defaultPassCounter = 1
	defaultDiceLength  = 6

	minUserPassLength = 8
	minGenPassLength  = 4
	maxGenPassLength  = 64
)

var (
	userName    string
	siteName    string
	passLength  uint8
	passCounter uint

	uppercase   bool
	lowercase   bool
	punctuation bool
	digits      bool

	isVerbose            bool
	doNotCopyToClipboard bool
	showPassword         bool
	behavior             string
)

func getEnvVars() {
	astolfoLength := os.Getenv("ASTOLFO_LENGTH")
	if astolfoLength != "" {
		astolfoLenUint8, err := strconv.ParseUint(astolfoLength, 0, 8)
		if err != nil {
			warn(err)
		} else {
			passLength = uint8(astolfoLenUint8)
		}
	}

	astolfoCounter := os.Getenv("ASTOLFO_COUNTER")
	if astolfoCounter != "" {
		astolfoCntUint, err := strconv.ParseUint(astolfoCounter, 0, 0)
		if err != nil {
			warn(err)
		} else {
			passCounter = uint(astolfoCntUint)
		}
	}

	astolfoMode = os.Getenv("ASTOLFO_MODE")
	if astolfoMode != "" {
		behavior = astolfoMode
	}
}

func initAndCheck() error {
	if userName == "" || siteName == "" {
		return errEmptyArg
	}
	if passLength < 4 || passLength > 64 {
		return errPassLength
	}
	if !uppercase && !lowercase && !punctuation && !digits {
		uppercase = true
		lowercase = true
		punctuation = true
		digits = true
	}

	behavior = strings.ToLower(behavior)
	switch behavior {
	case "hidecopy":
		return nil
	case "showcopy":
		showPassword = true
	case "showonly":
		showPassword = true
		doNotCopyToClipboard = true
	default:
		warn(fmt.Errorf("unknown behavior value; defaulting to \"hidecopy\""))
	}
	return nil
}

func inputPassword() ([]byte, error) {
	var passwd []byte
	var err error
	for {
		fmt.Fprint(os.Stderr, "Enter your master password: ")
		passwd, err = terminal.ReadPassword(int(syscall.Stdin))
		fmt.Fprint(os.Stderr, "\n")
		if err != nil {
			return nil, err
		}
		if utf8.RuneCountInString(string(passwd)) < minUserPassLength {
			fmt.Fprintf(os.Stderr, "Your password is too short! The minimum is %d letters\n", minUserPassLength)
			continue
		} else {
			return passwd, nil
		}
	}
}

func parseArg() error {
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Run `%s --help` for more information.\n", os.Args[0])
	}
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	pflag.CommandLine.SortFlags = false

	var displayVersion, displayHelp bool

	behaviorUsage := "set the preferred `behavior`; possible values:\n"
	behaviorUsage += "\"hidecopy\" (hide password, copy to clipboard)\n"
	behaviorUsage += "\"showcopy\" (show password, copy to clipboard)\n"
	behaviorUsage += "\"showonly\" (show password, don't copy to clipboard)\n"

	lengthUsage := "set the length of the generated password\n"
	lengthUsage += fmt.Sprintf("the `value` must range from %d to %d", minGenPassLength, maxGenPassLength)

	pflag.UintVarP(&passCounter, "counter", "c", defaultPassCounter, "set the counter `value`")
	pflag.Uint8VarP(&passLength, "length", "L", defaultPassLength, lengthUsage)
	pflag.BoolVarP(&digits, "digit", "d", false, "turn on numeric characters")
	pflag.BoolVarP(&lowercase, "lowercase", "l", false, "turn on lowercase letters")
	pflag.BoolVarP(&punctuation, "punctuation", "p", false, "turn on punctuation characters")
	pflag.BoolVarP(&uppercase, "uppercase", "U", false, "turn on uppercase letters")
	pflag.StringVarP(&behavior, "mode", "m", "hidecopy", behaviorUsage)
	pflag.BoolVarP(&isVerbose, "verbose", "v", false, "output more information")
	pflag.BoolVar(&displayVersion, "version", false, "display version information and exit")
	pflag.BoolVarP(&displayHelp, "help", "h", false, "display this help and exit")

	getEnvVars()

	pflag.Parse()

	if displayHelp {
		// If asking explicitly for --help, pflag.PrintDefaults() is
		// printed to stdout instead of stderr.
		pflag.CommandLine.SetOutput(os.Stdout)
		printHelp(os.Stdout)
		os.Exit(0)
	}
	if displayVersion {
		fmt.Println(version())
		os.Exit(0)
	}
	if n := pflag.NArg(); n < 2 || n > 2 {
		if n == 0 {
			printHelp(os.Stderr)
			os.Exit(0)
		}
		return errInsufficientArgs
	}
	userName = pflag.Arg(0)
	siteName = pflag.Arg(1)

	return nil
}

// printHelp prints the help usage. If the user explicitly passes the `--help`
// flag, it will write to the standard output. Otherwise, it will write
// to the standard error.
func printHelp(w io.Writer) {
	wr := bufio.NewWriter(w)
	wr.WriteString(fmt.Sprintf("%s - password generator\n", appName))

	wr.WriteString(fmt.Sprint("Usage:\n"))
	wr.WriteString(fmt.Sprintf("  %s [options] (<username> <sitename>)\n", os.Args[0]))
	wr.WriteString(fmt.Sprintf("  %s -h | --help\n", os.Args[0]))
	wr.WriteString(fmt.Sprintf("  %s --version\n\n", os.Args[0]))

	//fmt.Fprint(os.Stderr, "Example:\n")
	//fmt.Fprintf(os.Stderr, "  %s --length 4 --digit \"Cold Bishop\" \"example.com\"\n", appName)
	//fmt.Fprint(os.Stderr, "  Set the length to 4 and enable only numeric characters\n\n")

	wr.WriteString(fmt.Sprint("Options:\n"))
	wr.Flush()

	pflag.PrintDefaults()

	wr.WriteString(fmt.Sprint("\nIf none of the -dlpU options are passed, they are all enabled by default.\n"))
	wr.Flush()
}

func printUserParams() {
	fmt.Print("Params:\n")
	fmt.Printf("  Name: %s\n", userName)
	fmt.Printf("  Site: %s\n", siteName)
	fmt.Printf("  Password length: %d\n", passLength)
	fmt.Printf("  Password counter: %d\n", passCounter)

	fmt.Printf("  Uppercase: %t\n", uppercase)
	fmt.Printf("  Lowercase: %t\n", lowercase)
	fmt.Printf("  Punctuations: %t\n", punctuation)
	fmt.Printf("  Digits: %t\n\n", digits)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Nice little hack for returning from main() with a return value while
	// also calling all the deferred functions.
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	if err := parseArg(); err != nil {
		die(err, &exitCode, 1)
		return
	}
	if err := initAndCheck(); err != nil {
		die(err, &exitCode, 2)
		return
	}
	if isVerbose {
		printUserParams()
	}

	passwd, err := inputPassword()
	if err != nil {
		die(err, &exitCode, 3)
		return
	}

	generatedPassword, err := generatePassword(userName, siteName, passwd, passLength, passCounter, uppercase, lowercase, punctuation, digits)
	if err != nil {
		die(err, &exitCode, 4)
		return
	}

	if !doNotCopyToClipboard {
		err = clipboard.WriteAll(generatedPassword)
		if err != nil {
			warn(err)
			showPassword = true
		} else {
			fmt.Fprintln(os.Stderr, "Password generated and copied to the clipboard!")
		}
	}
	if showPassword {
		fmt.Fprint(os.Stderr, "Your password is: ")
		fmt.Println(generatedPassword)
	}
}
