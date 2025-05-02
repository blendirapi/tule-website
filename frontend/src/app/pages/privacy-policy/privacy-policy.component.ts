import { HeaderComponent } from '../../components/header/header.component';
import { FooterComponent } from '../../components/footer/footer.component';
import { Component } from '@angular/core';

@Component({
  selector: 'app-privacy-policy',
  standalone: true,
  imports: [HeaderComponent, FooterComponent],
  templateUrl: './privacy-policy.component.html',
})
export class PrivacyPolicyComponent {

}
