package gemini

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {

	payload := `{
  "candidates": [
    {
      "content": {
        "parts": [
          {
            "text": "Olá! Como posso ajudar você hoje?",
            "thoughtSignature": "EvAGCu0GARFNMg98NSHQOW/2iMzjV5n0Dgw0TZOmDMGpBbYxIMmdarRbnu3pF+zWKskqF6UjnkFuRYo7DvaoUjyI2Xf7qhUF9u7FGn22ND1GBgyVDJeSfHIqtib2NtlgA8Hc9K5CmVo6n5JtwsXss+mzwu3Aq4Z5iuLIp0/8QqG20LF3v9acg6f/euzW9SkgnPNYaK3fqALdxYJQjIIrZCWtcqU1BmRf0MoGIMh9WQcub3rUOgv0LPggl5xrGFBPzq0Ln/Yk6KSqwvLN0XWGs9tKICjpOZL/4kGLraxY9eHRdTD+wX2aIvkxTs3VwodrnPUNHmbSd6J/8hyhUFcYjiDrz3MovnUeNhmzUjYwD8lT9N+XqISKAOYqACEOYO+7EfR+TDUwH8xGhTBxUbKMnJJgaR5c9LVr9YPgvWI/oLFnNViJBL7Yh+GQgwREVR5DH/IKC1r8HBj4YjUUfdXzrtNXA5clzfOvRcALgynEQPa6hFQlEnflw7EBmupom1vzKo4mNpl5kj9wQAcdiAOLjMnps28YQ5cTbxPHaLWBAmIzsfnA1HoywSXHZ2ygndfnHVsslFuGM6vPUXUw4LQTmc53MrS2zq01AlJMs/WYCzUjyYlkKYBIdADT1bWgtChdBcRrvacKPoa0s3Q8Dic16Xb09VMP7yqU6HeybDP2Y/Bxv7PCb0FeGguRjT/Jh+q4wnRQPhUbRJqalThNEt9Xkj0I3YuFNjH0TCbwOLZatUQ4ZGrljt0Qy81PyYyNI3VAG1TpD+2/owcCCMsgO2Gh1NB5zs2QI37mqSpSQh9lDsyeXn/eR9cCJAwFOKgCyLImSnwOFZNddCd3Zm6UCx1tP4Ou/QXOev9uADrrFDHswY1sDcTQ/3pKg2lreoaxSNHG03N7GXqZ+lMzC+vhthT/dMyH2iYuPXckOZfQ1mMl/XXCTzxvT1hEwOOVllm8zs4kQqTKnIrIa/EyoQ/Ljdvp5J6Zw2Zj/UWEaahQAfIaU7PYbT/2+KIW2vbTp/HmYddcUo5FTanAsfYXj1hJCMX+Rf67kWwTsGyRYCcssHJJtoXKaAs1qIWaQ8z4HvP4XtkXVJvvzc+v3MFIEvZl4xfX32zMvBCFLzpQ2XnIH1d+nCiT9PwNjBg4mmQJtKIkvWvnt25HXZnbLeGJ1Mjf2BJ3C/qlQQ=="
          }
        ],
        "role": "model"
      },
      "finishReason": "STOP",
      "index": 0
    }
  ],
  "usageMetadata": {
    "promptTokenCount": 1,
    "candidatesTokenCount": 8,
    "totalTokenCount": 204,
    "promptTokensDetails": [
      {
        "modality": "TEXT",
        "tokenCount": 1
      }
    ],
    "thoughtsTokenCount": 195,
    "serviceTier": "standard"
  },
  "modelVersion": "gemini-3.6-flash",
  "responseId": "eBGfavudJabrz7IP-J7XYQ"
}`

	candidates, err := DecodePayload([]byte(payload))

	log.Print(err)
	assert.Nil(t, err)
	assert.NotNil(t, candidates)
}
