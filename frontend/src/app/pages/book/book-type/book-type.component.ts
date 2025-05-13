import { Component, Output, EventEmitter } from '@angular/core';

@Component({
	selector: 'app-book-type',
	standalone: true,
	imports: [],
	templateUrl: './book-type.component.html',
})
export class BookTypeComponent {
	@Output() optionSelected = new EventEmitter<string>();
	@Output() optionSelectedName = new EventEmitter<string>();
	@Output() price = new EventEmitter<number>();
	@Output() hasBath = new EventEmitter<string>();

	selectedOption: string | null = null;
	selectedOptionName: string = '';
	basePrice: number = 0;
	bath: string = '';
	shampooPrice: number = 3;

	selectOption(option: string, name: string, price: number) {
		this.selectedOption = option;
		this.selectedOptionName = name;
		this.basePrice = price;
	}

	toggleShampoo(event: Event, price: number) {
		const checkbox = event.target as HTMLInputElement;
		this.bath = checkbox.checked ? 'true' : 'false';
	}

	get totalPrice(): number {
		return this.basePrice + (this.bath ? this.shampooPrice : 0);
	}

	emitSelection() {
		if (!this.selectedOption) return;

		this.optionSelected.emit(this.selectedOption);
		this.optionSelectedName.emit(this.selectedOptionName);
		this.price.emit(this.totalPrice);
		this.hasBath.emit(this.bath);
	}
}
