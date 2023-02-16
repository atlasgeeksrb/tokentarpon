package main

import (
	"fmt"
	"net/http"
	"tokentarpon/tokenizer"

	"github.com/gin-gonic/gin"
)

var routerUrl = "localhost:8090"
var clientUrl = "localhost:3000"
var encryptionKey = "somefancylongkeyhere234*W&"

func main() {

	router := gin.Default()

	router.POST("/tokens", getTokens)
	router.PUT("/tokens", createTokens)
	router.OPTIONS("/tokens", preflight)

	router.POST("/tokens/values", getTokenValues)
	router.OPTIONS("/tokens/values", preflight)

	router.POST("/token", getToken)
	router.PUT("/token", createToken)
	router.DELETE("/token", deleteToken)
	router.OPTIONS("/token", preflight)

	router.POST("/token/value", getTokenValue)
	router.OPTIONS("/token/value", preflight)

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
	var tokenObj tokenizer.Token
	addHeaders(c)
	if err := c.BindJSON(&tokenObj); err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "Token is malformed"})
		return
	}
	domainUuid := tokenObj.DomainUuid
	uuid := tokenObj.Uuid
	//@todo, here we'd query something to figure out the name of the collection to use,
	tokenObj, err := tokenizer.GetToken(domainUuid, uuid)
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
	var tokenObj tokenizer.Token
	addHeaders(c)
	if err := c.BindJSON(&tokenObj); err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "Token is malformed"})
		return
	}
	domainUuid := tokenObj.DomainUuid
	uuid := tokenObj.Uuid
	//@todo, here we'd query something to figure out the name of the collection to use,
	tokenObj, err := tokenizer.GetToken(domainUuid, uuid)

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
	createdToken, mongoerr := tokenizer.CreateToken(tokenObj.DomainUuid, tokenObj.Value)
	if mongoerr != nil {
		errmsg := fmt.Sprint(mongoerr)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": errmsg})
	} else {
		c.IndentedJSON(http.StatusCreated, createdToken)
	}
}

func createTokens(c *gin.Context) {
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
	var domainUuid = tokens[0].DomainUuid
	createdTokens, errorTokens := tokenizer.CreateTokens(domainUuid, tokens)
	if len(errorTokens) > 0 {
		c.IndentedJSON(http.StatusInternalServerError, errorTokens)
	} else {
		c.IndentedJSON(http.StatusCreated, createdTokens)
	}
}

func deleteToken(c *gin.Context) {
	var tokenObj tokenizer.Token
	addHeaders(c)
	if err := c.BindJSON(&tokenObj); err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "Token record malformed"})
		return
	}

	//@todo, here we'd query something to figure out the name of the collection to use,
	// based on the domain Uuid
	// for now use the shared community store
	_, err := tokenizer.DeleteToken(tokenObj)
	if err != nil {
		errmsg := fmt.Sprint(err)
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": errmsg})
	} else {
		c.IndentedJSON(http.StatusCreated, gin.H{"message": "ok"})
	}
}

func getTokens(c *gin.Context) {
	var tokenObj tokenizer.Token
	addHeaders(c)
	if err := c.BindJSON(&tokenObj); err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "Token record malformed"})
		return
	}
	//@todo, here we'd query something to figure out the name of the collection to use,
	tokens, err := tokenizer.GetTokens(tokenObj.DomainUuid)
	if err != nil {
		errmsg := fmt.Sprint(err)
		c.IndentedJSON(http.StatusNotAcceptable, gin.H{"message": errmsg})
	} else {
		c.JSON(http.StatusOK, tokens)
	}
}

func getTokenValues(c *gin.Context) {
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
