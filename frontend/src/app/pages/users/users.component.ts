import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
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
	isEditing = false;

	constructor(private http: HttpClient) { }

	ngOnInit() {
		this.loadUsers();
	}

	loadUsers() {
		this.http.get<any[]>('/v0/api/users').subscribe((data) => {
			this.users = data;
		});
	}

	saveUser() {
		const apiUrl = this.isEditing ? `/v0/api/user/${this.form.user_id}` : '/v0/api/users';
		const method = this.isEditing ? 'put' : 'post';

		this.http[method](apiUrl, this.form).subscribe(() => {
			this.loadUsers();
			this.cancelEdit();
		});
	}

	editUser(user: any) {
		this.form = { ...user, password: '' };
		this.isEditing = true;
	}

	deleteUser(userId: number) {
		if (confirm('Are you sure you want to delete this user?')) {
			this.http.delete(`/v0/api/user/${userId}`).subscribe(() => this.loadUsers());
		}
	}

	cancelEdit() {
		this.form = { user_id: null, name: '', username: '', password: '' };
		this.isEditing = false;
	}
}
