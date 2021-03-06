import log4js from 'log4js';
import { IAuth } from '../../lib/auth';
import { IContext } from '../main';

export interface IManagePostBody {
    auth: IAuth;
    url: string;
    uuid: string;
}

const routeManagePost = async (ctx: IContext) => {
    const logger = log4js.getLogger('Server: /v3/manage/post');
    logger.level = ctx.logLevel;
    const body: IManagePostBody = ctx.request.body;
    const { pomment, userAuth } = ctx;
    if (!userAuth.auth(body.auth.time, body.auth.token)) {
        ctx.status = 403;
        return;
    }
    ctx.response.body = await pomment.getPost(body.url, body.uuid);
};

export default routeManagePost;
