import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import CryptoJS from 'crypto-js';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './login.component.html',
})
export class LoginComponent {
  formData = {
    username: '',
    password: ''
  };

  errorMessage = '';
  isSubmitting: boolean = false;

  constructor(private http: HttpClient, private router: Router) { }

  async onSubmit() {
	this.isSubmitting = true;
	try {
		const hash = CryptoJS.SHA256(CryptoJS.enc.Utf8.parse(this.formData.password)).toString();

		const payload = {
			username: this.formData.username,
			password: hash
		};

		this.http.post<any>('/v0/api/login', payload).subscribe({
			next: (response) => {
				localStorage.setItem('token', response.token);
				this.router.navigate(['/dashboard']);
			},
			error: (error) => {
				this.isSubmitting = false;
				this.errorMessage = error.error.message || 'Η σύνδεση απέτυχε. Παρακαλώ δοκιμάστε ξανά.';
			}
		});
	} catch (e) {
		this.errorMessage = 'Σφάλμα κρυπτογράφησης κωδικού.';
	}
}
}