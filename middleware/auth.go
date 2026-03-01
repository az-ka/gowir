package middleware

import (
	"net/http"
	// "gowir/internal/shared/response" // Akan digunakan nanti saat Auth sebenarnya diimplementasi
)

// RequireAuth adalah placeholder middleware untuk memastikan user sudah login.
// TODO: Implementasi verifikasi JWT / Session di sini.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Logika pengecekan token JWT akan diletakkan di sini.
		// Jika gagal:
		// response.Error(w, 401, "Unauthorized: Anda belum login")
		// return

		// Lanjut ke handler berikutnya jika sukses
		next.ServeHTTP(w, r)
	})
}

// RequireAdmin adalah placeholder middleware untuk memastikan user memiliki role "Admin".
// Harus diletakkan SETELAH RequireAuth.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Logika pengecekan role (misal dari claims JWT) akan diletakkan di sini.
		// Jika bukan admin:
		// response.Error(w, 403, "Forbidden: Anda tidak memiliki akses admin")
		// return

		// Lanjut ke handler berikutnya jika admin
		next.ServeHTTP(w, r)
	})
}
