import { Auth } from '../lib/auth';
import { PommentData } from '../core/main';
import { IConfig } from './config';

export interface IPommentContext {
    userConfig: IConfig;
    pomment: PommentData;
    logLevel: string;
    userAuth: Auth;
    userPath: string;
}
