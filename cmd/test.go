package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(fmt.Errorf("could not create new docker client: %w", err))
	}
	if err := client.FromEnv(cli); err != nil {
		panic(fmt.Errorf("could not configure docker client from environment: %w", err))
	}
	ctx := context.Background()
	imgPullOpts := types.ImagePullOptions{All: false}
	// imgPullOpts.RegistryAuth = "X2djcGlkZDpleUpoYkdjaU9pSlNVekkxTmlJc0ltdHBaQ0k2SWpjMk1tWmhOak0zWVdZNU5UTTFPVEJrWWpoaVlqaGhOak0yWW1ZeE1XUTBNell3WVdKak9UZ2lMQ0owZVhBaU9pSktWMVFpZlEuZXlKaGRXUWlPaUp5WldkcGMzUnllUzUyYW5CaGRHVnNMbTFsSWl3aVlYcHdJam9pTVRBeU1EazFPVFExTkRneU9USTBPRGN5TWpRd0lpd2laVzFoYVd3aU9pSXhNREU1TVRJeU56YzJPRGN4TFdOdmJYQjFkR1ZBWkdWMlpXeHZjR1Z5TG1kelpYSjJhV05sWVdOamIzVnVkQzVqYjIwaUxDSmxiV0ZwYkY5MlpYSnBabWxsWkNJNmRISjFaU3dpWlhod0lqb3hOVGd4T0RjNU1UazVMQ0puYjI5bmJHVWlPbnNpWTI5dGNIVjBaVjlsYm1kcGJtVWlPbnNpYVc1emRHRnVZMlZmWTNKbFlYUnBiMjVmZEdsdFpYTjBZVzF3SWpveE5UZ3hPRGN5Tmprd0xDSnBibk4wWVc1alpWOXBaQ0k2SWpRek9USTBOemd6TkRVNE5qVXlOVEV3TXpnaUxDSnBibk4wWVc1alpWOXVZVzFsSWpvaVoydGxMV3QxWW1Wc1pYUXRhVzFoWjJVdGMyVXRjSEpsWlcxd2RHbGliR1V0Ym05a1pTMDJOVE00TnprNE9TMWlOVFUzSWl3aWJHbGpaVzV6WlY5cFpDSTZXeUl4TURBeE1EQXpJaXdpTVRBd01UQXhNQ0lzSWpFMk5qY3pPVGN4TWpJek16WTFPRGMyTmlJc0lqWTRPREF3TkRFNU9EUXdPVFkxTkRBeE16SWlYU3dpY0hKdmFtVmpkRjlwWkNJNkluUmxjM1F0YTJsekxUazFOelZrTlRkaElpd2ljSEp2YW1WamRGOXVkVzFpWlhJaU9qRXdNVGt4TWpJM056WTROekVzSW5wdmJtVWlPaUpsZFhKdmNHVXRkMlZ6ZERFdFl5SjlmU3dpYVdGMElqb3hOVGd4T0RjMU5UazVMQ0pwYzNNaU9pSm9kSFJ3Y3pvdkwyRmpZMjkxYm5SekxtZHZiMmRzWlM1amIyMGlMQ0p6ZFdJaU9pSXhNREl3T1RVNU5EVTBPREk1TWpRNE56SXlOREFpZlEuTm9la2Vwck56R0FuSGVQX0N2dGxZakFfVExTM01qREJNUEJjaXdOcVV3UU9yMXRzTlI3ZXpzUE1yUUh3dnU3RzhGQnpCamJkbW9vbVJvOENOYi1YeWhpV3pUZ0NmRWdieF9FY1dHazVGUmw3clFPdjNZY1RHWDlnMzY1SlBxaXRiTEpfSmlNQzhweTZZY2p2R25mMnJkSW5zRzRXZWtibzB5RkdhOXdDN2RIOW1VdmJkZGlpZ1BnUHVKR2tMRmV5c29oQ3ZlUHBwMG9OM3Vuc3doeHNWSy1sel9MR2txSE8yUFNPbUJYSGpwM1BPQU5kbjVzYmh4WnZCckxsLWluazV5cldpRG8ycThKX0pINkZkc0ZiTm45LXlXbXJUVi0xTnVXTE0zNFY0ZHMzLTFxdjZaMXhzSzU3OFV2MWZPZTZ3Rm1sMHBXMGRhVnV0V3FWbWhTZThB"
	loginResp, err := cli.RegistryLogin(ctx, types.AuthConfig{Username: "vjftw", Password: "DnGpmRBW0dgLiIbPPtrI"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", loginResp.Status)
	regAuth := base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf(
			`{"username": "%s", "password": "%s"}`,
			"_gcpidd",
			"eyJhbGciOiJSUzI1NiIsImtpZCI6Ijc2MmZhNjM3YWY5NTM1OTBkYjhiYjhhNjM2YmYxMWQ0MzYwYWJjOTgiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJyZWdpc3RyeS52anBhdGVsLm1lIiwiYXpwIjoiMTAyMDk1OTQ1NDgyOTI0ODcyMjQwIiwiZW1haWwiOiIxMDE5MTIyNzc2ODcxLWNvbXB1dGVAZGV2ZWxvcGVyLmdzZXJ2aWNlYWNjb3VudC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiZXhwIjoxNTgxODc5MTk5LCJnb29nbGUiOnsiY29tcHV0ZV9lbmdpbmUiOnsiaW5zdGFuY2VfY3JlYXRpb25fdGltZXN0YW1wIjoxNTgxODcyNjkwLCJpbnN0YW5jZV9pZCI6IjQzOTI0NzgzNDU4NjUyNTEwMzgiLCJpbnN0YW5jZV9uYW1lIjoiZ2tlLWt1YmVsZXQtaW1hZ2Utc2UtcHJlZW1wdGlibGUtbm9kZS02NTM4Nzk4OS1iNTU3IiwibGljZW5zZV9pZCI6WyIxMDAxMDAzIiwiMTAwMTAxMCIsIjE2NjczOTcxMjIzMzY1ODc2NiIsIjY4ODAwNDE5ODQwOTY1NDAxMzIiXSwicHJvamVjdF9pZCI6InRlc3Qta2lzLTk1NzVkNTdhIiwicHJvamVjdF9udW1iZXIiOjEwMTkxMjI3NzY4NzEsInpvbmUiOiJldXJvcGUtd2VzdDEtYyJ9fSwiaWF0IjoxNTgxODc1NTk5LCJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJzdWIiOiIxMDIwOTU5NDU0ODI5MjQ4NzIyNDAifQ.NoekeprNzGAnHeP_CvtlYjA_TLS3MjDBMPBciwNqUwQOr1tsNR7ezsPMrQHwvu7G8FBzBjbdmoomRo8CNb-XyhiWzTgCfEgbx_EcWGk5FRl7rQOv3YcTGX9g365JPqitbLJ_JiMC8py6YcjvGnf2rdInsG4Wekbo0yFGa9wC7dH9mUvbddiigPgPuJGkLFeysohCvePpp0oN3unswhxsVK-lz_LGkqHO2PSOmBXHjp3POANdn5sbhxZvBrLl-ink5yrWiDo2q8J_JH6FdsFbNn9-yWmrTV-1NuWLM34V4ds3-1qv6Z1xsK578Uv1fOe6wFml0pW0daVutWqVmhSe8A",
		)),
	)
	imgPullOpts.RegistryAuth = regAuth
	out, err := cli.ImagePull(ctx, "registry.vjpatel.me/ghost:latest", imgPullOpts)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		os.Stdout.Write(scanner.Bytes())
	}
	out.Close()
}