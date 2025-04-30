import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import CryptoJS from 'crypto-js';
import { DashboardNavComponent } from '../../components/dashboard-nav/dashboard-nav.component';

@Component({
	selector: 'app-users',
	standalone: true,
	imports: [CommonModule, FormsModule, DashboardNavComponent],
	templateUrl: './users.component.html',
})
export class UsersComponent implements OnInit {
	users: any[] = [];
	form = { user_id: null, name: '', username: '', password: '' };
	isEditing: Boolean = false;
	isSubmitting: Boolean = false;

	constructor(private http: HttpClient) { }

	ngOnInit() {
		this.loadUsers();
	}

	loadUsers() {
		this.http.get<any[]>('/v0/api/users').subscribe((data) => {
			this.users = data;
			this.isSubmitting = false;
		});
	}

	async saveUser() {
		this.isSubmitting = true;
		try {
			const hash = CryptoJS.SHA256(CryptoJS.enc.Utf8.parse(this.form.password)).toString();

			const userPayload = {
				...this.form,
				password: hash
			};

			const apiUrl = this.isEditing ? `/v0/api/user/${this.form.user_id}` : `/v0/api/users`;
			const method = this.isEditing ? 'put' : 'post';

			this.http[method](apiUrl, userPayload).subscribe(() => {
				this.loadUsers();
				this.cancelEdit();
			});
		} catch (e) {
			console.error('Hashing failed:', e);
		}
	}

	editUser(user: any) {
		this.form = { ...user, password: '' };
		this.isEditing = true;
	}

	deleteUser(userId: number) {
		if (confirm('Είστε σίγουροι ότι θέλετε να διαγράψετε αυτόν τον χρήστη;')) {
			this.http.delete(`/v0/api/user/${userId}`).subscribe(() => this.loadUsers());
		}
	}

	cancelEdit() {
		this.form = { user_id: null, name: '', username: '', password: '' };
		this.isEditing = false;
	}
}
