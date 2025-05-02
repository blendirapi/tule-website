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

  selectOption(option: string, optionName: string) {
    this.optionSelected.emit(option);
    this.optionSelectedName.emit(optionName);
  }
}
