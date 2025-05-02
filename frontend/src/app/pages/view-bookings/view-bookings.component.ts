import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { format } from 'date-fns';
import { DashboardNavComponent } from '../../components/dashboard-nav/dashboard-nav.component';

interface Artist {
	user_id: number;
	name: string;
}

interface Service {
	service_id: number;
	name: string;
	time: string;
}

interface Booking {
	booking_id: number;
	firstname: string;
	lastname: string;
	email: string;
	phone: string;
	date: string;
	start_time: string;
	end_time: string;
	service_id: number;
	user_id: number;
	user_name: string;
	service_name: string;
	service_color: string;
}

@Component({
	selector: 'app-view-bookings',
	standalone: true,
	imports: [CommonModule, FormsModule, DashboardNavComponent],
	templateUrl: 'view-bookings.component.html',
})
export class DailyCalendarComponent implements OnInit {
	selectedBooking: Booking | null = null;
	newBooking: Partial<Booking> = {};
	newBookingModal: boolean = false;
	bookings: Booking[] = [];
	artists: Artist[] = [];
	services: Service[] = [];
	hours = Array.from({ length: 13 }, (_, i) => i + 9); // Here is the start time and also how many hours to show
	selectedDate = format(new Date(), 'yyyy-MM-dd');
	isSubmitting: Boolean = false;

	constructor(private http: HttpClient) { }

	ngOnInit() {
		this.loadBookings();
		this.loadArtists();
		this.loadServices();
	}

	loadBookings() {
		this.http.get<Booking[]>('/v0/api/get_bookings').subscribe({
			next: (data) => {
				this.bookings = data;
				this.isSubmitting = false;
			},
			error: (error) => {
				console.error('Failed to load bookings', error);
			}
		});
	}

	loadArtists() {
		this.http.get<Artist[]>('/v0/api/artists').subscribe({
			next: (data) => {
				this.artists = data;
			},
			error: (error) => {
				console.error('Failed to load artists', error);
			}
		});
	}

	loadServices() {
		this.http.get<Service[]>('/v0/api/services').subscribe({
			next: (data) => {
				this.services = data;
			},
			error: (error) => {
				console.error('Failed to load services', error);
			}
		});
	}

	getUserBookingsForDate(userId: number) {
		return this.bookings
			.filter(b => b.user_id === userId && this.getOnlyDate(b.date) === this.selectedDate)
			.sort((a, b) => a.start_time.localeCompare(b.start_time));
	}

	getOnlyDate(dateStr: string): string {
		return dateStr.split('T')[0];
	}

	getMinutesFromStart(time: string): number {
		const [hour, minute] = time.split(':').map(Number);
		return (hour - 9) * 60 + minute;
	}

	getDuration(start: string, end: string): number {
		const startParts = start.split(':').map(Number);
		const endParts = end.split(':').map(Number);
		return (endParts[0] * 60 + endParts[1]) - (startParts[0] * 60 + startParts[1]);
	}

	changeDate(days: number): void {
		const current = new Date(this.selectedDate);
		current.setDate(current.getDate() + days);
		this.selectedDate = current.toISOString().split('T')[0];
	}

	openBookingModal(booking: Booking) {
		this.selectedBooking = {...booking};

		if (this.selectedBooking?.date) {
			this.selectedBooking.date = format(new Date(this.selectedBooking.date), 'yyyy-MM-dd');
		}
	}

	closeModal() {
		this.selectedBooking = null;
	}

	openNewBookingModal() {
		this.newBookingModal = true;
	}
	
	closeNewBookingModal() {
		this.newBookingModal = false;
	}

	saveBookingChanges() {
		if (this.selectedBooking) {
			const updatedBooking = {
				artist: this.selectedBooking.user_id.toString(),
				service: this.selectedBooking.service_id.toString(),
				date: this.selectedBooking.date,
				time: this.selectedBooking.start_time + ' - ' + this.selectedBooking.end_time ,
				firstName: this.selectedBooking.firstname,
				lastName: this.selectedBooking.lastname,
				phone: this.selectedBooking.phone,
				email: this.selectedBooking.email
			};
	
			this.http.put(`/v0/api/bookings/${this.selectedBooking.booking_id}`, updatedBooking, {
				headers: { 'Content-Type': 'application/json' }
			})
			.subscribe({
				next: () => {
					this.closeModal();
					this.loadBookings();
				},
				error: (error) => {
					console.error('Failed to update booking:', error);
				}
			});
			
		}
	}	

	saveNewBooking() {
		if (this.newBooking) {
			const newBooking = {
				artist: this.newBooking.user_id?.toString(),
				service: this.newBooking.service_id?.toString(),
				date: this.newBooking.date,
				time: this.newBooking.start_time + ' - ' + this.newBooking.end_time ,
				firstName: this.newBooking.firstname,
				lastName: this.newBooking.lastname,
				phone: this.newBooking.phone,
				email: this.newBooking.email
			};
	
			this.http.post('/v0/api/add_dash_booking', newBooking, {
				headers: { 'Content-Type': 'application/json' }
			})
			.subscribe({
				next: () => {
					this.closeNewBookingModal();
					this.loadBookings();
				},
				error: (error) => {
					console.error('Failed to add booking:', error);
				}
			});
			
		}
	}	

	deleteBooking(id: number) {
		this.isSubmitting = true;
		this.http.delete(`/v0/api/bookings/${id}`)
			.subscribe({
				next: () => {
					this.closeModal();
					this.loadBookings();
				},
				error: (error) => {
					console.error('Failed to delete booking:', error);
				}
			});
	}

	adjustTime(field: 'start_time' | 'end_time', changeInMinutes: number) {
		if (!this.selectedBooking) return;

		const time = this.selectedBooking[field];
		if (!time) return;

		const [hours, minutes] = time.split(':').map(Number);
		let totalMinutes = hours * 60 + minutes + changeInMinutes;
		totalMinutes = Math.max(0, Math.min(1440, totalMinutes));

		const newHours = Math.floor(totalMinutes / 60);
		const newMinutes = totalMinutes % 60;

		this.selectedBooking[field] = `${String(newHours).padStart(2, '0')}:${String(newMinutes).padStart(2, '0')}`;
	}

	getUniqueUsers() {
		const users = this.bookings.map(b => ({ id: b.user_id, name: b.user_name }));
		const uniqueUsers = [...new Set(users.map(user => JSON.stringify(user)))].map(user => JSON.parse(user));
		uniqueUsers.sort((a, b) => a.id - b.id);
		return uniqueUsers;
	}	
}
