package db

func Float32ToFloat64(float32Slice []float32) []float64 {
    float64Slice := make([]float64, len(float32Slice))
    for i, val := range float32Slice {
        float64Slice[i] = float64(val)
    }
    return float64Slice
}
/**
cmd.Run() is for commands that execute and then exit:
The cmd.Run() function in Go's os/exec package is designed to:

Start the specified command.
Wait for the command to complete (exit).
Return any error that occurred during the command's execution or startup.
ollama serve is designed to run continuously as a server:
ollama serve is intended to start an Ollama server in the background.
It's a long-running process that listens for requests. It doesn't naturally "exit" after a short period.

Why cmd.Run() likely blocks or fails:

Blocking Behavior:  When you use cmd.Run() with ollama serve,
your Go program will likely block indefinitely, waiting for ollama serve to exit.
Since ollama serve is designed to run continuously, it might never exit naturally.
Your Go program will appear to hang or freeze.

Immediate Exit (Possible, but less likely):
In some scenarios, ollama serve might be designed to start in the background and then the initial ollama serve process might exit quickly after launching the background server process. In this less likely case,
cmd.Run() might return very quickly without an error,
but the desired server might or might not be properly started or accessible in the background.

Error if ollama serve fails to start immediately: If there's a problem launching ollama serve
(e.g., missing executable, port conflict, configuration issue), cmd.Run() will likely return an error.
However, even if it doesn't return an immediate error,
it's still probably not achieving what you want (a continuously running server).
**/
// func startOllamaServer() error {

// 	cmd := exec.Command("ollama", "serve")
// 	// time.Sleep(4 * time.Second)
// 	err := cmd.Run()
// 	if err != nil {
// 		log.Println("Error:", err)
// 		return err
// 	}
// 	log.Println("ollama server is running")
// 	// log.Println(string(out))
// 	return nil
// }
