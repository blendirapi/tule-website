import { Injectable } from "@angular/core";

@Injectable({ providedIn: 'root' })
export class AuthService {
	isLoggedIn(): boolean {
		const token = localStorage.getItem('auth');
		return !!token;
	}

	logout(): void {
		localStorage.removeItem('auth');
	}
}
