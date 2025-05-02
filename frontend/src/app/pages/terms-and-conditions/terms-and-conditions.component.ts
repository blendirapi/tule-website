import { HeaderComponent } from '../../components/header/header.component';
import { FooterComponent } from '../../components/footer/footer.component';
import { Component } from '@angular/core';

@Component({
  selector: 'app-terms-and-conditions',
  standalone: true,
  imports: [HeaderComponent, FooterComponent],
  templateUrl: './terms-and-conditions.component.html',
})
export class TermsAndConditionsComponent {

}
