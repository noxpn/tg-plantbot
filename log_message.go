package main

import "fmt"

//log to console json -> pretty
func logMessage(update Update)  {
	fmt.Printf(" msg_id: %d \n offset: %d\n user_id: %d\n user_name: %s\n",
		update.Message.Id,
		update.UpdateID,
		update.Message.Chat.Id,
		update.Message.Chat.Name)
	if len(update.Message.Text) > 0 {
		fmt.Printf(" text: %s\n", update.Message.Text)
	} else {
		fmt.Printf(" text: \"empty\"\n", )
	}
	if len(update.Message.Photo) > 0 {
		logPhoto(update)
	}
	if !update.Message.Document.isEmpty() {
		logDocument(update)
	}

	fmt.Printf("\n")
	
}

func logPhoto(update Update) {
	fmt.Printf(" photos: %d\n", len(update.Message.Photo))
	for _, v := range update.Message.Photo {
		fmt.Printf("  uniq_id: %s\n   id: %s", v.Unique_id, v.Id)	
	}
}

func logDocument(update Update) {
	fmt.Printf("  uniq_id: %s\n  file: %s\n  size: %d\n", 
	update.Message.Document.Unique_id,
	update.Message.Document.Name,
	update.Message.Document.Size)	
}