package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	port       = getEnv("PORT", "4001")
	resendKey  = os.Getenv("RESEND_API_KEY")
	toEmail    = getEnv("TO_EMAIL", "kalkowski123@gmail.com")
	fromEmail  = getEnv("FROM_EMAIL", "onboarding@resend.dev")
)

type ContactForm struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

type ResendEmail struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
	ReplyTo string   `json:"reply_to"`
}

type APIResponse struct {
	Success bool   `json:"success,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	http.HandleFunc("/api/contact", handleContact)
	http.HandleFunc("/", handleStatic)

	log.Printf("🚀 Server running at http://localhost:%s", port)
	log.Printf("📧 Emails will be sent to: %s", toEmail)
	if resendKey == "" {
		log.Println("⚠️  Warning: RESEND_API_KEY not set. Email sending will fail.")
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}

	// Prevent directory traversal
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Try to serve from current directory
	filePath := "." + cleanPath

	// Check if file exists
	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		http.NotFound(w, r)
		return
	}

	// Set content type
	ext := filepath.Ext(filePath)
	contentType := getContentType(ext)
	w.Header().Set("Content-Type", contentType)

	// Serve file
	data, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func getContentType(ext string) string {
	types := map[string]string{
		".html": "text/html; charset=utf-8",
		".css":  "text/css; charset=utf-8",
		".js":   "application/javascript",
		".json": "application/json",
		".png":  "image/png",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".gif":  "image/gif",
		".svg":  "image/svg+xml",
		".ico":  "image/x-icon",
	}
	if ct, ok := types[ext]; ok {
		return ct
	}
	return "application/octet-stream"
}

func handleContact(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Error: "Method not allowed"})
		return
	}

	// Parse request body
	var form ContactForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Error: "Invalid JSON"})
		return
	}

	// Validate required fields
	if form.Name == "" || form.Email == "" || form.Message == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Error: "Wymagane pola: imię, email, wiadomość"})
		return
	}

	// Build email HTML
	phoneDisplay := form.Phone
	if phoneDisplay == "" {
		phoneDisplay = "Nie podano"
	}

	html := fmt.Sprintf(`
		<h2>Nowa wiadomość ze strony Balkonowa Ochrona</h2>
		<table style="border-collapse: collapse; width: 100%%; max-width: 600px;">
			<tr>
				<td style="padding: 10px; border: 1px solid #ddd; font-weight: bold;">Imię i nazwisko:</td>
				<td style="padding: 10px; border: 1px solid #ddd;">%s</td>
			</tr>
			<tr>
				<td style="padding: 10px; border: 1px solid #ddd; font-weight: bold;">Email:</td>
				<td style="padding: 10px; border: 1px solid #ddd;"><a href="mailto:%s">%s</a></td>
			</tr>
			<tr>
				<td style="padding: 10px; border: 1px solid #ddd; font-weight: bold;">Telefon:</td>
				<td style="padding: 10px; border: 1px solid #ddd;">%s</td>
			</tr>
			<tr>
				<td style="padding: 10px; border: 1px solid #ddd; font-weight: bold;">Wiadomość:</td>
				<td style="padding: 10px; border: 1px solid #ddd;">%s</td>
			</tr>
		</table>
	`, form.Name, form.Email, form.Email, phoneDisplay, strings.ReplaceAll(form.Message, "\n", "<br>"))

	// Send email via Resend
	email := ResendEmail{
		From:    fromEmail,
		To:      []string{toEmail},
		Subject: fmt.Sprintf("Nowa wiadomość od %s - Balkonowa Ochrona", form.Name),
		HTML:    html,
		ReplyTo: form.Email,
	}

	emailJSON, _ := json.Marshal(email)

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(emailJSON))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{Error: "Błąd wysyłania wiadomości"})
		return
	}

	req.Header.Set("Authorization", "Bearer "+resendKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{Error: "Błąd wysyłania wiadomości"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Resend error: %s", string(body))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{Error: "Błąd wysyłania wiadomości"})
		return
	}

	log.Printf("Email sent successfully to %s from %s", toEmail, form.Email)
	json.NewEncoder(w).Encode(APIResponse{Success: true, Message: "Wiadomość została wysłana"})
}