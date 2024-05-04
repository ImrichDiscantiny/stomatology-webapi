/*
 * Appointment List Api
 *
 * Stomatology Appointment List management for Web-In-Cloud system
 *
 * API version: 1.0.0
 * Contact: xdiscantiny@stuba.sk
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package stomatology_al

import (
    "github.com/gin-gonic/gin"
)

func AddRoutes(engine *gin.Engine) {
  group := engine.Group("/api")
  
  {
    api := newStomatologyAppointmentListAPI()
    api.addRoutes(group)
  }
  
}