import { Component, Output, EventEmitter } from '@angular/core';

@Component({
  selector: 'app-book-type',
  standalone: true,
  imports: [],
  templateUrl: './book-type.component.html',
})
export class BookTypeComponent {
  @Output() optionSelected = new EventEmitter<string>();
  @Output() back = new EventEmitter<void>();

  selectOption(option: string) {
    this.optionSelected.emit(option);
  }

  goBack(): void {
    this.back.emit();
}   
}
