package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {}

func (cli *CLI) printUsage() {
	fmt.Println("Usage: gchain <command> <flag>...")
	fmt.Println(" * createblockchain -addr <address>")
	fmt.Println("     : Create a blockchain and send the genesis block reward to <address>.")
	fmt.Println(" * createwallet")
	fmt.Println("     : Generate a new key-pair and save it into the wallet.")
	fmt.Println(" * getbalance -addr <address>")
	fmt.Println("     : Get the balance of <address>.")
	fmt.Println(" * listaddr")
	fmt.Println("     : List all the addresses from the wallet.")
	fmt.Println(" * printchain")
	fmt.Println("     : Print all the blocks of the blockchain.")
	fmt.Println(" * reindexutxo")
	fmt.Println("     : Rebuild the UTXO set.")
	fmt.Println(" * send -from <from> -to <to> -amount <amount> -mine")
	fmt.Println("     : Send <amount> of coins from <from> address to <to> address.")
	fmt.Println("       Mine on the same node, when -mine is set.")
	fmt.Println(" * startnode -miner <miner>")
	fmt.Println("     : Start a node with ID specified in NODE_ID env. var.")
	fmt.Println("       -miner enables mining and send the block reward to <miner> address.")
}

func (cli *CLI) printUsageAndExit() {
	cli.printUsage()
	os.Exit(1)
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsageAndExit()
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Println("NODE_ID env. var. is not set.")
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "createblockchain":
		err = cli.handleCreateBlockchain(nodeID, os.Args[2:])
	case "createwallet":
		err = cli.handleCreateWallet(nodeID, os.Args[2:])
	case "getbalance":
		err = cli.handleGetBalance(nodeID, os.Args[2:])
	case "listaddr":
		err = cli.handleListAddresses(nodeID, os.Args[2:])
	case "printchain":
		err = cli.handlePrintChain(nodeID, os.Args[2:])
	case "reindexutxo":
		err = cli.handleReindexUTXO(nodeID, os.Args[2:])
	case "send":
		err = cli.handleSend(nodeID, os.Args[2:])
	case "startnode":
		err = cli.handleStartNode(nodeID, os.Args[2:])
	default:
		cli.printUsageAndExit()
	}

	if err != nil {
		log.Panic(err)
	}
}

func (cli *CLI) handleCreateBlockchain(nodeID string, flags []string) error {
	cmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	addr := cmd.String("addr", "", "The address to send the genesis block reward to")

	if err := cmd.Parse(flags); err != nil {
		return err
	}

	if *addr == "" {
		cmd.Usage()
		os.Exit(1)
	}

	createBlockchain(nodeID, *addr)
	return nil
}

func (cli *CLI) handleCreateWallet(nodeID string, flags []string) error {
	cmd := flag.NewFlagSet("createwallet", flag.ExitOnError)

	if err := cmd.Parse(flags); err != nil {
		return err
	}

	createWallet(nodeID)
	return nil
}

func (cli *CLI) handleGetBalance(nodeID string, flags []string) error {
	cmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	addr := cmd.String("addr", "", "The address to get balance for")

	if err := cmd.Parse(flags); err != nil {
		return err
	}

	if *addr == "" {
		cmd.Usage()
		os.Exit(1)
	}

	getBalance(nodeID, *addr)
	return nil
}

func (cli *CLI) handleListAddresses(nodeID string, flags []string) error {
	cmd := flag.NewFlagSet("listaddr", flag.ExitOnError)

	if err := cmd.Parse(flags); err != nil {
		return err
	}

	listAddresses(nodeID)
	return nil
}

func (cli *CLI) handlePrintChain(nodeID string, flags []string) error {
	cmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	if err := cmd.Parse(flags); err != nil {
		return err
	}

	printChain(nodeID)
	return nil
}

func (cli *CLI) handleReindexUTXO(nodeID string, flags []string) error {
	cmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)

	if err := cmd.Parse(flags); err != nil {
		return err
	}

	reindexUTXO(nodeID)
	return nil
}

func (cli *CLI) handleSend(nodeID string, flags []string) error {
	cmd := flag.NewFlagSet("send", flag.ExitOnError)
	from := cmd.String("from", "", "The source address to send coins from")
	to := cmd.String("to", "", "The destination address to send coins to")
	amount := cmd.Int("amount", 0, "The amount of coins to send")
	mine := cmd.Bool("mine", false, "The mine flag to decide whether mining immediately on the same node.")

	if err := cmd.Parse(flags); err != nil {
		return err
	}

	if *from == "" || *to == "" || *amount <= 0 {
		cmd.Usage()
		os.Exit(1)
	}

	send(nodeID, *from, *to, *amount, *mine)
	return nil
}

func (cli *CLI) handleStartNode(nodeID string, flags []string) error {
	cmd := flag.NewFlagSet("startnode", flag.ExitOnError)
	miner := cmd.String("miner", "", "The miner address to enables mining and send the block reward to")

	if err := cmd.Parse(flags); err != nil {
		return err
	}

	startNode(nodeID, *miner)
	return nil
}

func NewCLI() *CLI {
	return &CLI{}
}
