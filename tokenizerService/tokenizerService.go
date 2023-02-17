package main

import (
	"fmt"
	"net/http"
	"strconv"
	"tokentarpon/tokenizer"

	"github.com/gin-gonic/gin"
)

var routerUrl = "localhost:8090"
var clientUrl = "localhost:3000"
var pageSize int64 = 100

//var encryptionKey = "somefancylongkeyhere234*W&"

func main() {

	router := gin.Default()

	router.GET("/tokens/:domainId", getTokens)
	router.PUT("/tokens/:domainId", createTokens)
	router.OPTIONS("/tokens/:domainId", preflight)

	router.POST("/tokens/:domainId/values", getTokenValues)
	router.OPTIONS("/tokens/:domainId/values", preflight)

	router.GET("/tokens/:domainId/:id", getToken)
	router.PUT("/tokens/:domainId/:id", createToken)
	router.DELETE("/tokens/:domainId/:id", deleteToken)
	router.OPTIONS("/tokens/:domainId/:id", preflight)

	router.GET("/tokens/:domainId/:id/value", getTokenValue)
	router.OPTIONS("/tokens/:domainId/:id/value", preflight)

	router.Run(routerUrl)

}

func addOptionsHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", clientUrl) //"*"
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers,x-auth-token,content-type")
	c.Header("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,DELETE,PUT")
}

func addHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", clientUrl) //"*"
	c.Header("Access-Control-Expose-Headers", "x-auth-token")
}

func preflight(c *gin.Context) {
	addOptionsHeaders(c)
	c.JSON(http.StatusOK, struct{}{})
}

func getToken(c *gin.Context) {
	domainUuid := c.Param("domainId")
	tokenId := c.Param("id")
	addHeaders(c)

	//@todo, here we'd query something to figure out the name of the collection to use,
	tokenObj, err := tokenizer.GetToken(domainUuid, tokenId)
	if err != nil {
		errmsg := fmt.Sprint(err)
		// if err == datastore.ErrNotFound {
		// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "token not found"})
		// } else if err == datastore.ErrServerError || err == datastore.ErrDatastoreError {
		// 	c.IndentedJSON(http.StatusNotAcceptable, gin.H{"message": errmsg})
		// } else {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": errmsg})
		//}
	} else {
		c.IndentedJSON(http.StatusOK, tokenObj)
	}
}

func getTokenValue(c *gin.Context) {
	domainUuid := c.Param("domainId")
	tokenId := c.Param("id")

	addHeaders(c)

	//@todo, here we'd query something to figure out the name of the collection to use,
	tokenObj, err := tokenizer.GetToken(domainUuid, tokenId)

	if err != nil {
		errmsg := fmt.Sprint(err)
		// if err == datastore.ErrNotFound {
		// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "token not found"})
		// } else if err == datastore.ErrServerError || err == datastore.ErrDatastoreError {
		// 	c.IndentedJSON(http.StatusNotAcceptable, gin.H{"message": errmsg})
		// } else {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": errmsg})
		//}
	} else {
		c.IndentedJSON(http.StatusOK, tokenObj.Value)
	}
}

func createToken(c *gin.Context) {
	domainUuid := c.Param("domainId")

	var tokenObj tokenizer.Token
	addHeaders(c)
	if err := c.BindJSON(&tokenObj); err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "Token record malformed"})
		return
	}

	//@todo, here we'd query something to figure out the name of the collection to use,
	// based on the domain Uuid
	// e.g. tokenizer.CollectionName = "mycollection"
	// for now use the shared community store
	createdToken, dataerr := tokenizer.CreateToken(domainUuid, tokenObj.Value)
	if dataerr != nil {
		errmsg := fmt.Sprint(dataerr)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": errmsg})
	} else {
		c.IndentedJSON(http.StatusCreated, createdToken)
	}
}

func createTokens(c *gin.Context) {
	domainUuid := c.Param("domainId")
	var tokens []tokenizer.Token
	addHeaders(c)
	if err := c.BindJSON(&tokens); err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "Token record malformed"})
		return
	}

	if len(tokens) == 0 {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "No tokens were provided"})
		return
	}

	//@todo, here we'd query something to figure out the name of the collection to use,
	// based on the domain Uuid
	// for now use the shared community store
	createdTokens, errorTokens := tokenizer.CreateTokens(domainUuid, tokens)
	if len(errorTokens) > 0 {
		c.IndentedJSON(http.StatusInternalServerError, errorTokens)
	} else {
		c.IndentedJSON(http.StatusCreated, createdTokens)
	}
}

func deleteToken(c *gin.Context) {
	domainUuid := c.Param("domainId")
	tokenId := c.Param("id")
	addHeaders(c)

	//@todo, here we'd query something to figure out the name of the collection to use,
	// based on the domain Uuid
	// for now use the shared community store
	_, err := tokenizer.DeleteToken(domainUuid, tokenId)
	if err != nil {
		errmsg := fmt.Sprint(err)
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": errmsg})
	} else {
		c.IndentedJSON(http.StatusCreated, gin.H{"message": "ok"})
	}
}

func getTokens(c *gin.Context) {
	domainUuid := c.Param("domainId")
	start, limit := getPageParams(c)
	addHeaders(c)

	//@todo, here we'd query something to figure out the name of the collection to use,
	tokens, err := tokenizer.GetTokens(domainUuid, start, limit)
	if err != nil {
		errmsg := fmt.Sprint(err)
		c.IndentedJSON(http.StatusNotAcceptable, gin.H{"message": errmsg})
	} else {
		c.JSON(http.StatusOK, tokens)
	}
}

func getTokenValues(c *gin.Context) {
	//domainUuid := c.Param("domainId")
	//start, limit := getPageParams(c)
	var tokenQuery tokenizer.TokenQuery
	addHeaders(c)
	if err := c.BindJSON(&tokenQuery); err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "Token request malformed"})
		return
	}

	//@todo here replace community with the user's collection
	tokenValues, err := tokenizer.GetTokenValues(tokenQuery)
	if err != nil {
		errmsg := fmt.Sprint(err)
		c.IndentedJSON(http.StatusNotAcceptable, gin.H{"message": errmsg})
	} else {
		c.JSON(http.StatusOK, tokenValues)
	}
}

func getPageParams(c *gin.Context) (int64, int64) {
	// get start, limit from queryparams
	var start int64
	var limit int64
	if startparam, ok := c.GetQuery("start"); ok {
		i, err := strconv.ParseInt(startparam, 10, 64)
		if err == nil {
			if i < 0 {
				i = 0
			} else {
				start = i
			}
		} else {
			i = 0
		}
	}
	if limitparam, ok := c.GetQuery("limit"); ok {
		i, err := strconv.ParseInt(limitparam, 10, 64)
		if err == nil {
			limit = i
		}
	} else {
		limit = pageSize
	}
	return start, limit
}
