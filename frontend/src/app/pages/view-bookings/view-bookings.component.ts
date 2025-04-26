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
}

@Component({
	selector: 'app-view-bookings',
	standalone: true,
	imports: [CommonModule, FormsModule, DashboardNavComponent],
	templateUrl: 'view-bookings.component.html',
})
export class DailyCalendarComponent implements OnInit {
	selectedBooking: Booking | null = null;
	bookingCopy: Booking | null = null;
	bookings: Booking[] = [];
	artists: Artist[] = [];
	hours = Array.from({ length: 14 }, (_, i) => i + 8);
	selectedDate = format(new Date(), 'yyyy-MM-dd');

	constructor(private http: HttpClient) { }

	ngOnInit() {
		this.loadBookings();
		this.loadArtists();
	}

	loadBookings() {
		this.http.get<Booking[]>('/v0/api/get_bookings').subscribe({
			next: (data) => {
				this.bookings = data;
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
		return (hour - 8) * 60 + minute;
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
		this.selectedBooking = booking;
		this.bookingCopy = { ...booking };

		if (this.bookingCopy?.date) {
			this.bookingCopy.date = format(new Date(this.bookingCopy.date), 'yyyy-MM-dd');
		}
	}

	closeModal() {
		this.selectedBooking = null;
		this.bookingCopy = null;
	}

	saveBookingChanges() {
		if (this.selectedBooking && this.bookingCopy) {
			Object.assign(this.selectedBooking, this.bookingCopy);
	
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
				next: (response) => {
					this.closeModal();
					this.loadBookings();
				},
				error: (error) => {
					console.error('Failed to update booking:', error);
				}
			});
			
		}
	}	

	deleteBooking(id: number) {
		this.http.delete(`/v0/api/bookings/${id}`)
			.subscribe({
				next: (response) => {
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
		return [...new Set(users.map(user => JSON.stringify(user)))].map(user => JSON.parse(user));
	}
}
