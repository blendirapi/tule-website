import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { DashboardNavComponent } from '../../components/dashboard-nav/dashboard-nav.component';

type WeekDay = 'monday' | 'tuesday' | 'wednesday' | 'thursday' | 'friday' | 'saturday' | 'sunday';

interface WeekDayItem {
	key: WeekDay;
	label: string;
}

@Component({
	selector: 'app-users',
	standalone: true,
	imports: [CommonModule, FormsModule, DashboardNavComponent],
	templateUrl: './users.component.html',
})
export class UsersComponent implements OnInit {
	users: any[] = [];
	isEditing = false;
	isSubmitting = false;

	weekDays: WeekDayItem[] = [
		{ key: 'monday', label: 'Δευτέρα' },
		{ key: 'tuesday', label: 'Τρίτη' },
		{ key: 'wednesday', label: 'Τετάρτη' },
		{ key: 'thursday', label: 'Πέμπτη' },
		{ key: 'friday', label: 'Παρασκευή' },
		{ key: 'saturday', label: 'Σάββατο' },
		{ key: 'sunday', label: 'Κυριακή' }
	];

	form: ReturnType<typeof this.getEmptyForm> = this.getEmptyForm();

	constructor(private http: HttpClient) { }

	ngOnInit() {
		this.loadUsers();
	}

	getEmptyForm() {
		const createDaySlots = () => [
			{ start: '00:00', end: '00:00' },
			{ start: '00:00', end: '00:00' }
		];

		return {
			userId: null,
			name: '',
			username: '',
			email: '',
			password: '',
			isVisible: false,
			workingHours: {
				monday: createDaySlots(),
				tuesday: createDaySlots(),
				wednesday: createDaySlots(),
				thursday: createDaySlots(),
				friday: createDaySlots(),
				saturday: createDaySlots(),
				sunday: createDaySlots()
			}
		};
	}

	loadUsers() {
		this.http.get<any[]>('/v0/api/users').subscribe((data) => {
			this.users = data;
			this.isSubmitting = false;
			this.cancelEdit();
		});
	}

	async saveUser() {
		this.isSubmitting = true;
		try {
			const userPayload = {
				userId: this.form.userId,
				name: this.form.name,
				username: this.form.username,
				email: this.form.email,
				password: this.form.password,
				isVisible: this.form.isVisible,
				workingHours: this.form.workingHours
			};

			const apiUrl = this.isEditing ? `/v0/api/user/${this.form.userId}` : `/v0/api/users`;
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
		const emptyForm = this.getEmptyForm();
		const mergedWorkingHours: any = {};

		for (const day of this.weekDays) {
			const key = day.key;
			const userDaySlots = user.workingHours?.[key] ?? [];

			mergedWorkingHours[key] = [
				userDaySlots[0] ?? { start: '00:00', end: '00:00' },
				userDaySlots[1] ?? { start: '00:00', end: '00:00' },
			];
		}

		this.form = {
			...emptyForm,
			...structuredClone(user),
			workingHours: mergedWorkingHours,
			password: ''
		};

		this.isEditing = true;
	}

	deleteUser(userId: number) {
		if (confirm('Είστε σίγουροι ότι θέλετε να διαγράψετε αυτόν τον χρήστη;')) {
			this.http.delete(`/v0/api/user/${userId}`).subscribe(() => this.loadUsers());
		}
	}

	cancelEdit() {
		this.isEditing = false;
		this.form = this.getEmptyForm();
	}
}
