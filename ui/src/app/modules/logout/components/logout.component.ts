import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AppComponent } from '../../../app.component';
import { BACKEND_AUTH_URL } from '../../../app.constants';

@Component({
    selector: 'app-logout',
    templateUrl: './logout.component.html',
    styleUrls: ['./logout.component.scss']
})
export class LogoutComponent implements OnInit {
    public form: any = {};
    public LOGIN_STATUS = "wait";
    ngOnInit() {
        this.handleLogout();
        this.LOGIN_STATUS = "wait";
    }

    constructor(private router: Router, private http: HttpClient, private appComponent: AppComponent) { }

    public handleLogout() {
        let logoutURL = BACKEND_AUTH_URL + 'logout';
        const logoutOptions = {
            withCredentials: true
        };
        this.http.post(logoutURL, JSON.stringify({}), logoutOptions).subscribe((response) => {
            this.appComponent.IS_LOGEDIN = false;
        }, (err) => {
            console.log("Error", err);
        }
        );
        this.router.navigateByUrl("./login");
    }
}