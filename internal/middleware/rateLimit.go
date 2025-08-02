package middleware

// type RateLimiter struct {
// 	// Implement rate limiter logic here
// 	globalRateLimiter         *limiter.Limiter
// 	publicAPIRateLimiter      *limiter.Limiter
// 	userPrivateAPIRateLimiter *limiter.Limiter
// }

// func NewRateLimiter() *RateLimiter {
// 	rateLimit := &RateLimiter{
// 		globalRateLimiter:         rateLimiter("100-S"),
// 		publicAPIRateLimiter:      rateLimiter("80-S"),
// 		userPrivateAPIRateLimiter: rateLimiter("50-S"),
// 	}

// 	return rateLimit
// }

// func rateLimiter(interval string) *limiter.Limiter {
// 	store, err := redisStore.NewStoreWithOptions(global.Rdb, limiter.StoreOptions{
// 		Prefix:          "rate-limiter",
// 		MaxRetry:        3,
// 		CleanUpInterval: time.Hour,
// 	})

// 	if err != nil {
// 		return nil
// 	}

// 	rate, err := limiter.NewRateFromFormatted(interval) // 5-S, 10-M
// 	if err != nil {
// 		panic(err)
// 	}

// 	instance := limiter.New(store, rate)

// 	return instance
// }

// // global limiter
// func (rl *RateLimiter) GlobalRateLimiter() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		key := "global" // unit
// 		log.Println("global --->")

// 		limitContext, err := rl.globalRateLimiter.Get(c, key)
// 		if err != nil {
// 			fmt.Println("Failed to check rate limit Global", err)
// 			c.Next()
// 			return
// 		}

// 		if limitContext.Reached {
// 			log.Printf("Rate limit breached Global %s", key)
// 			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
// 				"error": "Rate limit breached GLOBAL, try later",
// 			})

// 			return
// 		}

// 		c.Next()
// 	}
// }

// // public api limiter
// func (rl *RateLimiter) PublicAPIRateLimiter() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		urlPath := c.Request.URL.Path
// 		rateLimitPath := rl.filterLimitUrlPath(urlPath)

// 		if rateLimitPath != nil {
// 			log.Println("Client IP --->", c.ClientIP())
// 			key := fmt.Sprintf("%s-%s", "111-222-333-444", urlPath)
// 			limitContext, err := rateLimitPath.Get(c, key)

// 			if err != nil {
// 				fmt.Println("Failed to check rate limit", err)
// 				c.Next()
// 				return
// 			}

// 			if limitContext.Reached {
// 				log.Printf("Rate limit breached %s", key)
// 				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
// 					"error": "Rate limit breached, try later",
// 				})
// 			}
// 		}

// 		c.Next()
// 	}
// }

// // private api limiter
// func (rl *RateLimiter) UserPrivateAPIRateLimiter() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		urlPath := c.Request.URL.Path // url request
// 		rateLimitPath := rl.filterLimitUrlPath(urlPath)

// 		if rateLimitPath != nil {
// 			userId := 1001 // uuid
// 			key := fmt.Sprintf("%s-%s", userId, urlPath)
// 			limitContext, err := rateLimitPath.Get(c, key)

// 			if err != nil {
// 				fmt.Println("Failed to check rate limit", err)
// 				c.Next()
// 				return
// 			}

// 			if limitContext.Reached {
// 				log.Printf("Rate limit breached %s", key)
// 				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
// 					"error": "Rate limit breached, try later",
// 				})
// 			}
// 		}

// 		c.Next()
// 	}
// }

// func (rl *RateLimiter) filterLimitUrlPath(urlPath string) *limiter.Limiter {
// 	if urlPath == "/v1/2025/user/login" || urlPath == "/ping/80" {
// 		return rl.publicAPIRateLimiter
// 	} else if urlPath == "/v1/2025/user/info" || urlPath == "/ping/50" {
// 		return rl.userPrivateAPIRateLimiter
// 	} else {
// 		return rl.globalRateLimiter
// 	}
// }
