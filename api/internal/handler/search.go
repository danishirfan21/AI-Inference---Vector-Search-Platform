package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-repo/ai-platform/api/internal/service"
)

import (
	"encoding/json"
	"time"
)

type SearchHandler struct {
	natsService   *service.NATSService
	cacheService  *service.CacheService
	milvusService *service.MilvusService
}

func NewSearchHandler(nats *service.NATSService, cache *service.CacheService, milvus *service.MilvusService) *SearchHandler {
	return &SearchHandler{natsService: nats, cacheService: cache, milvusService: milvus}
}

func (h *SearchHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query is required"})
		return
	}

	// Try cache
	if cached, err := h.cacheService.Get(c.Request.Context(), "search:"+query); err == nil {
		var results []interface{}
		if err := json.Unmarshal([]byte(cached), &results); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"query":   query,
				"results": results,
				"cached":  true,
			})
			return
		}
	}

	embedding, err := h.natsService.RequestEmbedding(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get embedding"})
		return
	}

	ids, err := h.milvusService.Search(c.Request.Context(), "text_embeddings", embedding)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search Milvus"})
		return
	}

	// Cache results
	resData, _ := json.Marshal(ids)
	h.cacheService.Set(c.Request.Context(), "search:"+query, resData, 10*time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"query":   query,
		"results": ids,
	})
}
