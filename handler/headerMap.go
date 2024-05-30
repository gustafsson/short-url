package handler

import "net/http"

func headersToMap(headers http.Header) map[string]interface{} {
    headerMap := make(map[string]interface{})
    for key, values := range headers {
        // Convert the slice of strings to a single string, or keep it as a slice
        // Depending on your needs, you might want to join the values or keep them as a slice
        if len(values) == 1 {
            headerMap[key] = values[0] // Single value as string
        } else {
            headerMap[key] = values // Multiple values as slice of strings
        }
    }
    return headerMap
}
