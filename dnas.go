package main

func main() {
	CLIRun(func(options *Options) {
		Monitor(options)
	})
}
