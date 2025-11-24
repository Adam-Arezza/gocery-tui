package services

import (
	"encoding/json"
	"fmt"
	"net/http"
    "io"
    "bytes"
	"github.com/Adam-Arezza/gocery-tui/config"
	"github.com/Adam-Arezza/gocery-tui/internal/types"
)

func GetGroceryItems(serverCfg *config.ServerConfig) ([]types.GroceryItem, error){
    url := "http://" + serverCfg.Host + ":" + serverCfg.Port + "/grocery_items"        
    resp, err := http.Get(url)
    if err != nil {
        return nil,err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API error: %d", resp.StatusCode)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var groceryItems []types.GroceryItem
    err = json.Unmarshal(body, &groceryItems)
    if err != nil {
        return nil, err
    }
    return groceryItems, nil
}

func UpdateGrocery(serverCfg *config.ServerConfig, items []types.PurchaseRequestItem) error{
    url := "http://" + serverCfg.Host + ":" + serverCfg.Port + "/grocery_items"
	jsonData, err := json.Marshal(items)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
        return err
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
    return nil
}
