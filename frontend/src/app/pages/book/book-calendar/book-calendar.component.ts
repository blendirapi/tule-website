import { Component, Input, OnInit, Output, EventEmitter } from '@angular/core';
import { CommonModule, registerLocaleData } from '@angular/common';
import { DateAdapter } from '@angular/material/core';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { provideNativeDateAdapter } from '@angular/material/core';
import { LOCALE_ID } from '@angular/core';
import localeEl from '@angular/common/locales/el';
import { HttpClient, HttpParams } from '@angular/common/http';

registerLocaleData(localeEl);

interface DateAvailability {
	date: Date;
	availability: 'high' | 'medium' | 'low' | 'none';
}

interface ApiAvailabilityResponse {
	success: boolean;
	classes: { [date: string]: string };
}

@Component({
	selector: 'app-book-calendar',
	imports: [CommonModule, MatDatepickerModule],
	providers: [
		{ provide: LOCALE_ID, useValue: 'el-GR' },
		provideNativeDateAdapter()
	],
	templateUrl: './book-calendar.component.html',
})
export class BookCalendarComponent implements OnInit {
	minDate: Date = new Date();
	maxDate: Date = new Date();
	selectedDate: Date | null = null;
	dateAvailabilities: DateAvailability[] = [];

	@Input() serviceType: string = '';
	@Input() selectedArtist: string | null = null;
	@Output() dateSelected = new EventEmitter<Date>();

	constructor(
		private dateAdapter: DateAdapter<Date>,
		private http: HttpClient
	) {
		this.dateAdapter.setLocale('el-GR');
	}

	ngOnInit(): void {
		const today = new Date();
		this.minDate = new Date(today.getFullYear(), today.getMonth(), today.getDate());
		this.maxDate = new Date(today.getFullYear(), today.getMonth() + 3, 0);
		this.fetchAvailabilityData();
	}

	fetchAvailabilityData(): void {
		if (!this.selectedArtist) {
			console.error('Artist ID is required');
			return;
		}

		let params = new HttpParams()
			.set('artist_id', this.selectedArtist);

		this.http.get<ApiAvailabilityResponse>('/v0/api/availability', { params })
			.subscribe({
				next: (response) => {
					if (response.success && response.classes) {
						this.processAvailabilityData(response.classes);
					}
				},
				error: (error) => {
					console.error('Error fetching availability data:', error);
				}
			});
	}

	processAvailabilityData(classes: { [date: string]: string }): void {
		const today = new Date();
		const endDate = new Date(today.getFullYear(), today.getMonth() + 3, 0);

		let currentDate = new Date(today);
		while (currentDate <= endDate) {
			const dateStr = this.formatDate(currentDate);
			const availabilityClass = classes[dateStr] || 'none';

			this.dateAvailabilities.push({
				date: new Date(currentDate),
				availability: availabilityClass as 'high' | 'medium' | 'low' | 'none',
			});

			currentDate.setDate(currentDate.getDate() + 1);
		}
	}

	private formatDate(date: Date): string {
		const year = date.getFullYear();
		const month = (date.getMonth() + 1).toString().padStart(2, '0');
		const day = date.getDate().toString().padStart(2, '0');
		return `${year}-${month}-${day}`;
	}

	dateFilter = (date: Date | null): boolean => {
		if (!date) return false;
		return this.getDateAvailability(date) !== 'none';
	};

	getDateAvailability(date: Date): 'high' | 'medium' | 'low' | 'none' {
		const found = this.dateAvailabilities.find(
			(d) =>
				d.date.getDate() === date.getDate() &&
				d.date.getMonth() === date.getMonth() &&
				d.date.getFullYear() === date.getFullYear()
		);
		return found ? found.availability : 'none';
	};

	getDateClass = (date: Date): string => {
		const availability = this.getDateAvailability(date);
		return `available-${availability}`;
	};

	onDateSelected(date: Date | null): void {
		this.selectedDate = date;
	}

	proceedWithDate(): void {
		if (this.selectedDate) {
			this.dateSelected.emit(this.selectedDate);
		}
	}
}