import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ElementRef, ViewChild, HostListener } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpClient, HttpParams } from '@angular/common/http';
import { format, startOfWeek, endOfWeek, addWeeks } from 'date-fns';
import { el } from 'date-fns/locale';
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
	bath: boolean;
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
	hours = Array.from({ length: 12 }, (_, i) => i + 9); // Here is the start time and also how many hours to show
	selectedDate = format(new Date(), 'yyyy-MM-dd');
	isSubmitting: Boolean = false;
	viewMode: 'daily' | 'weekly' = 'daily';
	currentWeekOffset = 0;
	suggestedClients: Partial<Booking>[] = [];
	clientSearchTimeout: any;
	@ViewChild('phoneInputWrapper') phoneInputWrapper!: ElementRef;
	showSuggestions: boolean = false;
	activeBookingType: 'new' | 'selected' | null = null;

	constructor(private http: HttpClient) { }

	ngOnInit() {
		const storedMode = localStorage.getItem('calendarViewMode');

		if (storedMode === 'daily' || storedMode === 'weekly') {
			this.viewMode = storedMode;
		} else {
			this.viewMode = 'daily';
		}

		if (this.viewMode === 'weekly') {
			this.currentWeekOffset = 0;
		}

		this.loadBookings();
	}


	loadBookings() {
		this.loadArtists();
		this.loadServices();

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
		this.http.get<Artist[]>('/v0/api/dash_artists').subscribe({
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
		this.selectedBooking = { ...booking };

		if (this.selectedBooking?.date) {
			this.selectedBooking.date = format(new Date(this.selectedBooking.date), 'yyyy-MM-dd');
		}
	}

	closeModal() {
		this.selectedBooking = null;
	}

	openNewBookingModal() {
		this.newBookingModal = true;
		this.newBooking = {
			date: this.selectedDate,
			start_time: '00:00',
			end_time: '00:00'
		};
	}

	closeNewBookingModal() {
		this.newBookingModal = false;
		this.newBooking = {};
	}

	saveBookingChanges() {
		if (this.selectedBooking) {
			this.isSubmitting = true;
			const updatedBooking = {
				artist: this.selectedBooking.user_id.toString(),
				service: this.selectedBooking.service_id.toString(),
				date: this.selectedBooking.date,
				time: this.selectedBooking.start_time + ' - ' + this.selectedBooking.end_time,
				firstName: this.selectedBooking.firstname,
				lastName: this.selectedBooking.lastname,
				phone: this.selectedBooking.phone,
				email: this.selectedBooking.email,
				bath: this.newBooking.bath ? "true" : "false"
			};

			this.http.put(`/v0/api/bookings/${this.selectedBooking.booking_id}`, updatedBooking, {
				headers: { 'Content-Type': 'application/json' }
			})
				.subscribe({
					next: () => {
						this.closeModal();
						this.loadBookings();
						this.isSubmitting = false;
					},
					error: (error) => {
						console.error('Failed to update booking:', error);
					}
				});

		}
	}

	saveNewBooking() {
		if (this.newBooking) {
			this.isSubmitting = true;
			const newBooking = {
				artist: this.newBooking.user_id?.toString(),
				service: this.newBooking.service_id?.toString(),
				date: this.newBooking.date,
				time: this.newBooking.start_time + ' - ' + this.newBooking.end_time,
				firstName: this.newBooking.firstname,
				lastName: this.newBooking.lastname,
				phone: this.newBooking.phone,
				email: this.newBooking.email,
				bath: this.newBooking.bath ? "true" : "false"
			};

			this.http.post('/v0/api/add_dash_booking', newBooking, {
				headers: { 'Content-Type': 'application/json' }
			})
				.subscribe({
					next: () => {
						this.closeNewBookingModal();
						this.loadBookings();
						this.isSubmitting = false;
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

	adjustTime(
		field: 'start_time' | 'end_time',
		changeInMinutes: number,
		target: 'selected' | 'new') {
		const booking = target === 'selected' ? this.selectedBooking : this.newBooking;
		if (!booking) return;

		const time = booking[field];
		if (!time) return;

		const [hours, minutes] = time.split(':').map(Number);
		let totalMinutes = hours * 60 + minutes + changeInMinutes;
		totalMinutes = Math.max(0, Math.min(1440, totalMinutes));

		const newHours = Math.floor(totalMinutes / 60);
		const newMinutes = totalMinutes % 60;

		booking[field] = `${String(newHours).padStart(2, '0')}:${String(newMinutes).padStart(2, '0')}`;
	}

	getUniqueUsers() {
		const users = this.bookings.map(b => ({ id: b.user_id, name: b.user_name }));
		const uniqueUsers = [...new Set(users.map(user => JSON.stringify(user)))].map(user => JSON.parse(user));
		uniqueUsers.sort((a, b) => a.id - b.id);
		return uniqueUsers;
	}

	getFormattedSelectedDate(): string {
		const baseDate = new Date(this.selectedDate);

		if (this.viewMode === 'weekly') {
			// Apply week offset (in days)
			const offsetDate = addWeeks(baseDate, this.currentWeekOffset);

			// Week starts on Monday and ends on Saturday
			const startOfWeekDate = startOfWeek(offsetDate, { weekStartsOn: 1 });
			const endOfWeekDate = endOfWeek(offsetDate, { weekStartsOn: 0 });

			const formattedStart = format(startOfWeekDate, 'EEEE, dd MMMM yyyy', { locale: el });
			const formattedEnd = format(endOfWeekDate, 'EEEE, dd MMMM yyyy', { locale: el });

			return `${formattedStart} - ${formattedEnd}`;
		} else {
			return format(baseDate, 'EEEE, dd MMMM yyyy', { locale: el });
		}
	}

	goToToday() {
		this.loadBookings();
		this.selectedDate = format(new Date(), 'yyyy-MM-dd');
	}

	onServiceChange() {
		if (this.newBooking.service_id && this.newBooking.start_time) {
			const serviceId = +this.newBooking.service_id;
			const selectedService = this.services.find(service => service.service_id === serviceId);

			if (selectedService) {
				// Parse start time "HH:mm" string into hours and minutes
				const [startHours, startMinutes] = this.newBooking.start_time.split(':').map(Number);

				const [serviceHours, serviceMinutes] = selectedService.time.split(':').map(Number);
				const serviceTimeInMinutes = (serviceHours * 60) + serviceMinutes;

				const startDate = new Date();

				startDate.setHours(startHours, startMinutes, 0, 0);
				startDate.setMinutes(startDate.getMinutes() + serviceTimeInMinutes);

				const endHours = String(startDate.getHours()).padStart(2, '0');
				const endMinutes = String(startDate.getMinutes()).padStart(2, '0');

				this.newBooking.end_time = `${endHours}:${endMinutes}`;
			}
		}
	}

	getStartOfWeek(date: Date): Date {
		const d = new Date(date);
		const day = d.getDay(); // 0 (Sun) to 6 (Sat)
		const diff = d.getDate() - ((day + 6) % 7); // Adjust so Monday = 0
		d.setDate(diff);
		d.setHours(0, 0, 0, 0);
		return d;
	}

	// Returns bookings for a user on a specific day
	getUserBookingsForDateWeekly(userId: number, date: Date): Booking[] {
		return this.bookings.filter(b =>
			b.user_id === userId &&
			new Date(b.date).toDateString() === date.toDateString()
		);
	}

	getGreekDayName(date: Date): string {
		const days = [
			'Κυριακή',
			'Δευτέρα',
			'Τρίτη',
			'Τετάρτη',
			'Πέμπτη',
			'Παρασκευή',
			'Σάββατο'
		];
		return days[date.getDay()];
	}

	switchView(mode: 'daily' | 'weekly') {
		this.viewMode = mode;
		localStorage.setItem('calendarViewMode', mode)
		if (mode === 'weekly') {
			this.currentWeekOffset = 0;
		}
	}

	changeWeek(offset: number) {
		this.currentWeekOffset += offset;
	}

	goToCurrentWeek() {
		this.currentWeekOffset = 0;
	}

	// Modify getWeekDays to respect currentWeekOffset
	getWeekDays(): Date[] {
		const baseDate = new Date(this.selectedDate);
		const startOfWeek = this.getStartOfWeek(baseDate);
		startOfWeek.setDate(startOfWeek.getDate() + this.currentWeekOffset * 7);
		return Array.from({ length: 6 }, (_, i) => {
			const d = new Date(startOfWeek);
			d.setDate(d.getDate() + i);
			return d;
		});
	}

	disableDay(event: Event, date: string, userId: number) {
		const isChecked = (event.target as HTMLInputElement).checked;

		this.isSubmitting = true;
		const disableData = {
			status: isChecked,
			date: date,
			userId: userId,
		};

		this.http.post('/v0/api/disable_day', disableData, {
			headers: { 'Content-Type': 'application/json' }
		})
			.subscribe({
				next: () => {
					this.loadBookings();
					this.isSubmitting = false;
				},
				error: (error) => {
					console.error('Failed to add booking:', error);
				}
			});
	}

	hasFullDayBooking(date: string, userId: number): boolean {
		const userBookings = this.bookings.filter(b =>
			b.user_id === userId &&
			this.getOnlyDate(b.date) === date &&
			b.start_time === '09:00' &&
			b.end_time === '21:00'
		);

		return userBookings.length > 0 ? true : false;
	}

	activateSuggestions(type: 'new' | 'selected') {
		this.activeBookingType = type;
		this.showSuggestions = true;
	}

	onPhoneInput(phone: string, type: 'new' | 'selected') {
		this.activeBookingType = type;
		clearTimeout(this.clientSearchTimeout);

		if (!phone || phone.length < 5) {
			this.suggestedClients = [];
			return;
		}

		// Inside your method:
		this.clientSearchTimeout = setTimeout(() => {
			const params = new HttpParams().set('phone', phone);

			this.http.get<Partial<Booking>[]>('/v0/api/get_client_data', { params }).subscribe({
				next: (data) => {
					this.suggestedClients = data ?? [];
				},
				error: (error) => {
					console.error('Client lookup failed:', error);
				}
			});
		}, 300); // 300ms debounce

	}

	selectClient(client: Partial<Booking>) {
		const target = this.activeBookingType === 'selected' ? this.selectedBooking : this.newBooking;

		if (target) {
			target.firstname = client.firstname || '';
			target.lastname = client.lastname || '';
			target.email = client.email || '';
			target.phone = client.phone || '';
		}

		this.suggestedClients = [];
		this.showSuggestions = false;
	}
}
