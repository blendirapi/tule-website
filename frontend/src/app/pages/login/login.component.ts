import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';

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
		const payload = {
			username: this.formData.username,
			password: this.formData.password
		};

		this.http.post<any>('/v0/api/login', payload).subscribe({
			next: (response) => {
				localStorage.setItem('auth', response.token);
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

  forgotPassword() {
	this.http.post('/v0/api/forgot_password', this.formData.username, {
	headers: { 'Content-Type': 'application/json' }
	})
	.subscribe({
		next: () => {
			console.log('Password reset email sent successfully.');
		},
		error: (error) => {
			console.error('Failed to add booking:', error);
		}
	});
  }
}
