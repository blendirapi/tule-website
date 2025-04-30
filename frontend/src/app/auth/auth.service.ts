import { Injectable } from "@angular/core";

@Injectable({ providedIn: 'root' })
export class AuthService {
	isLoggedIn(): boolean {
		const token = localStorage.getItem('token');
		// You could also add token expiration check here
		return !!token;
	}

	logout(): void {
		localStorage.removeItem('token');
	}
}
