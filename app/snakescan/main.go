package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dsxg666/cubery"
	"github.com/dsxg666/snakecoin/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:8545", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}
	client := pb.NewRpcClient(conn)

	app := cubery.Default()
	app.Static("/static", "./static")
	app.LoadHTMLGlob("templates/**/*")
	app.NoRoute(func(c *cubery.Context) {
		c.HTML(http.StatusNotFound, "main/404.html", nil)
	})
	app.GET("/", func(c *cubery.Context) {
		r, _ := client.GetAllBlock(context.Background(), &pb.GetAllBlockReq{})
		r2, _ := client.GetAllTx(context.Background(), &pb.GetAllTxReq{})
		c.HTML(http.StatusOK, "main/index.html", cubery.H{
			"blocks": r.Bs[:6],
			"txs":    r2.Txs[:6],
		})
	})
	app.GET("/txs", func(c *cubery.Context) {
		pageNum := c.Query("pageNum")
		currentPage, _ := strconv.Atoi(pageNum)
		r, _ := client.GetAllTx(context.Background(), &pb.GetAllTxReq{})
		totalRecords := len(r.Txs)
		totalPages := totalRecords / 9
		if totalRecords%9 != 0 {
			totalPages++
		}
		var start, end, left, right int
		if currentPage == 1 {
			left = 1
			right = 2
		} else if currentPage == totalPages {
			left = totalPages - 1
			right = totalPages
		} else {
			left = currentPage - 1
			right = currentPage + 1
		}
		if totalRecords%9 != 0 && currentPage == totalPages {
			start = (currentPage - 1) * 9
			end = start + totalRecords%9
		} else {
			start = (currentPage - 1) * 9
			end = start + 9
		}
		c.HTML(http.StatusOK, "main/txs.html", cubery.H{
			"txs":         r.Txs[start:end],
			"totalPage":   totalPages,
			"currentPage": currentPage,
			"left":        left,
			"right":       right,
		})
	})
	app.GET("/blocks", func(c *cubery.Context) {
		pageNum := c.Query("pageNum")
		currentPage, _ := strconv.Atoi(pageNum)
		r, _ := client.GetAllBlock(context.Background(), &pb.GetAllBlockReq{})
		totalRecords := len(r.Bs)
		totalPages := totalRecords / 9
		if totalRecords%9 != 0 {
			totalPages++
		}
		var start, end, left, right int
		if currentPage == 1 {
			left = 1
			right = 2
		} else if currentPage == totalPages {
			left = totalPages - 1
			right = totalPages
		} else {
			left = currentPage - 1
			right = currentPage + 1
		}
		if totalRecords%9 != 0 && currentPage == totalPages {
			start = (currentPage - 1) * 9
			end = start + totalRecords%9
		} else {
			start = (currentPage - 1) * 9
			end = start + 9
		}
		c.HTML(http.StatusOK, "main/blocks.html", cubery.H{
			"blocks":      r.Bs[start:end],
			"totalPage":   totalPages,
			"currentPage": currentPage,
			"left":        left,
			"right":       right,
		})
	})
	app.GET("/tx/:hash", func(c *cubery.Context) {
		r, _ := client.GetInfoByTxHash(context.Background(), &pb.GetInfoByTxHashReq{Hash: c.Param("hash")})
		if r == nil {
			c.HTML(http.StatusOK, "main/unfind.html", nil)
		} else {
			c.HTML(http.StatusOK, "main/tx.html", cubery.H{
				"txHash":   r.GetTxHash()[2:],
				"blockNum": r.GetBlock(),
				"time":     r.GetTime(),
				"from":     r.GetFrom(),
				"to":       r.GetTo(),
				"amount":   r.GetAmount(),
			})
		}
	})
	app.GET("/txPool", func(c *cubery.Context) {
		r, _ := client.GetTxPool(context.Background(), &pb.GetTxPoolReq{})
		if r == nil {
			c.HTML(http.StatusOK, "main/pool.html", nil)
		} else {
			c.HTML(http.StatusOK, "main/pool.html", cubery.H{
				"txs": r.GetTxs(),
			})
		}
	})
	app.GET("/newAccount", func(c *cubery.Context) {
		c.HTML(http.StatusOK, "main/newAccount.html", nil)
	})
	app.POST("/newAccount", func(c *cubery.Context) {
		r, _ := client.NewAccount(context.Background(), &pb.NewAccountReq{Password: c.PostForm("password")})
		c.HTML(http.StatusOK, "main/success.html", cubery.H{"address": r.GetAccount()})
	})
	app.GET("/newTx", func(c *cubery.Context) {
		c.HTML(http.StatusOK, "main/newTx.html", nil)
	})
	app.POST("/newTx", func(c *cubery.Context) {
		r, _ := client.NewTx(context.Background(), &pb.NewTxReq{
			From:     c.PostForm("from"),
			To:       c.PostForm("to"),
			Amount:   c.PostForm("amount"),
			Password: c.PostForm("password"),
		})
		if r.GetState() == "0" {
			c.JSON(http.StatusOK, cubery.H{
				"state": "0",
			})
		} else if r.GetState() == "1" {
			c.JSON(http.StatusOK, cubery.H{
				"state": "1",
			})
		} else if r.GetState() == "2" {
			c.JSON(http.StatusOK, cubery.H{
				"state": "2",
			})
		} else if r.GetState() == "3" {
			c.JSON(http.StatusOK, cubery.H{
				"state": "3",
			})
		} else if r.GetState() == "4" {
			c.JSON(http.StatusOK, cubery.H{
				"state": "4",
			})
		} else if r.GetState() == "5" {
			c.JSON(http.StatusOK, cubery.H{
				"state": "5",
			})
		}
	})
	app.GET("/block/:num", func(c *cubery.Context) {
		num := c.Param("num")
		if num == "byHash" {
			hash := c.Query("hash")
			r, _ := client.GetInfoByBlockHash(context.Background(), &pb.GetInfoByBlockHashReq{Hash: hash})
			if r == nil {
				c.HTML(http.StatusOK, "main/unfind.html", nil)
			} else {
				c.HTML(http.StatusOK, "main/block.html", cubery.H{
					"num":            num,
					"number":         r.GetNumber(),
					"nonce":          r.GetNonce(),
					"time":           r.GetTime(),
					"txs":            r.GetTxs(),
					"reward":         r.GetReward(),
					"difficulty":     r.GetDifficulty(),
					"coinbase":       r.GetCoinbase(),
					"blockHash":      r.GetBlockHash()[2:],
					"prevBlockHash":  r.GetPrevBlockHash()[2:],
					"stateTreeRoot":  r.GetStateTreeRoot()[2:],
					"merkleTreeRoot": r.GetMerkleTreeRoot()[2:],
				})
			}
		} else {
			r, _ := client.GetInfoByBlockNum(context.Background(), &pb.GetInfoByBlockNumReq{Num: num})
			if r == nil {
				c.HTML(http.StatusOK, "main/unfind.html", nil)
			} else {
				c.HTML(http.StatusOK, "main/block.html", cubery.H{
					"num":            num,
					"number":         r.GetNumber(),
					"nonce":          r.GetNonce(),
					"time":           r.GetTime(),
					"txs":            r.GetTxs(),
					"reward":         r.GetReward(),
					"difficulty":     r.GetDifficulty(),
					"coinbase":       r.GetCoinbase(),
					"blockHash":      r.GetBlockHash()[2:],
					"prevBlockHash":  r.GetPrevBlockHash()[2:],
					"stateTreeRoot":  r.GetStateTreeRoot()[2:],
					"merkleTreeRoot": r.GetMerkleTreeRoot()[2:],
				})
			}
		}
	})

	app.GET("/mine", func(c *cubery.Context) {
		c.HTML(http.StatusOK, "main/mine.html", nil)
	})
	app.POST("/mine", func(c *cubery.Context) {
		r, _ := client.Mine(context.Background(), &pb.MineReq{})
		c.JSON(http.StatusOK, cubery.H{"nonce": r.GetNonce()})
	})
	app.POST("/newBlock", func(c *cubery.Context) {
		r, _ := client.NewBlock(context.Background(), &pb.NewBlockReq{
			Nonce: c.PostForm("nonce"),
			Miner: c.PostForm("miner"),
		})
		c.JSON(http.StatusOK, cubery.H{"state": r.GetState()})
	})
	app.GET("/getBalance", func(c *cubery.Context) {
		c.HTML(http.StatusOK, "main/balance.html", nil)
	})
	app.POST("/getBalance", func(c *cubery.Context) {
		r, _ := client.GetBalance(context.Background(), &pb.GetBalanceReq{Addr: c.PostForm("addr")})
		if r == nil {
			c.JSON(http.StatusOK, cubery.H{
				"balance": "unExist",
			})
		} else {
			c.JSON(http.StatusOK, cubery.H{
				"balance": r.GetBalance(),
			})
		}
	})
	err = app.Run(":8080")
	if err != nil {
		fmt.Println("app.Run have something wrong! Err: ", err)
	}
}
