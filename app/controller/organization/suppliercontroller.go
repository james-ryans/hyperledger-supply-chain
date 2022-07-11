package controller

// func GetMeAsSupplier(c *gin.Context) {
// 	gateway, err := fabric.Gateway()
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": err.Error(),
// 			"data":    nil,
// 		})
//
// 		return
// 	}
// 	defer gateway.Close()
//
// 	certificate, err := identity.CertificateFromPEM(gateway.Identity().Credentials())
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": err.Error(),
// 			"data":    nil,
// 		})
//
// 		return
// 	}
//
// 	certificateJSON, err := json.Marshal(certificate)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": err.Error(),
// 			"data":    nil,
// 		})
//
// 		return
// 	}
//
// 	var test gin.H
// 	json.Unmarshal(certificate.Extensions[len(certificate.Extensions)-1].Value, &test)
// 	fmt.Printf("%s\n", test)
//
// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": nil,
// 		"data":    string(certificateJSON),
// 	})
// }

// func InitSupplier(c *gin.Context) {
// 	gateway, err := fabric.Gateway()
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": err.Error(),
// 			"data":    nil,
// 		})
//
// 		return
// 	}
// 	defer gateway.Close()
//
// 	network := gateway.GetNetwork(os.Getenv("FABRIC_CHANNEL_NAME"))
// 	contract := network.GetContract(os.Getenv("FABRIC_CHAINCODE_NAME"))
//
// 	id := "28899ce0-b79d-497b-ae78-fd3b896e0429"
//
// 	_, err = contract.SubmitTransaction("SupplierContract:CreateSupplier", id, "Supplier 0", "North Sumatra", "Medan", "Medan Kota", "20212", "Jl. Thamrin", "081234567890", "supplier0@hehe.com", "1234.1", "4321.1")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": err.Error(),
// 			"data":    nil,
// 		})
//
// 		return
// 	}
//
// 	supplierJSON, err := contract.EvaluateTransaction("SupplierContract:ReadSupplier", id)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": err.Error(),
// 			"data":    nil,
// 		})
//
// 		return
// 	}
// 	var supplier model.Supplier
// 	err = json.Unmarshal(supplierJSON, &supplier)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": err.Error(),
// 			"data":    nil,
// 		})
//
// 		return
// 	}
//
// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": nil,
// 		"data":    supplier,
// 	})
// }
