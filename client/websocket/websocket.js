/**
 * 输出信息到textarea
 * @param content
 */
function appendToDialog(content) {
    $(".log").append(content + "\r\n");
}

const protoVersion = 1
// size
const rawHeaderLen = 16;
// offset
const packetOffset = 0;
const headerOffset = 4;
const verOffset = 6;
const opOffset = 8;
const seqOffset = 12;
// op
const OpHeartbeatReply = 3
const OpAuth = 7
const OpAuthReply = 8
const OpBatchMsg = 9

class WebsocketOp {
    WsClient;

    constructor() {
        this.textEncoder = new TextEncoder();
        this.textDecoder = new TextDecoder();
    }

    // 连接
    connect() {
        let url = "ws://127.0.0.1:12572";// ws://是web socket协议,发送到websocket服务器的16779端口
        this.WsClient = new WebSocket(url);
        this.WsClient.binaryType = 'arraybuffer';

        /**
         * 当WebSocket对象的readyState状态变为OPEN时会触发该事件。
         * 该事件表明websocket连接成功并开始发送数据。
         * event属于Event对象，https://developer.mozilla.org/en-US/docs/Web/API/Event
         */
        this.WsClient.onopen = (event) => {
            appendToDialog("连接成功...");
            console.log(event);
            this.auth()
        };
        /**
         * 当有消息到达客户端的时候该事件会触发
         * event属于MessageEvent对象，https://developer.mozilla.org/en-US/docs/Web/API/MessageEvent
         */
        this.WsClient.onmessage = (event) => {
            console.log(event.data)
            let data = event.data;
            let dataView = new DataView(data, 0);
            let packetLen = dataView.getInt32(packetOffset);
            let headerLen = dataView.getInt16(headerOffset);
            let ver = dataView.getInt16(verOffset);
            let op = dataView.getInt32(opOffset);
            let seq = dataView.getInt32(seqOffset);
            // appendToDialog('获得消息：' + event.data);
            // console.log("receiveHeader: packetLen=" + packetLen, "headerLen=" + headerLen, "ver=" + ver, "op=" + op, "seq=" + seq);

            switch (op) {
                case OpAuthReply:
                    appendToDialog('授权成功...');
                    // send a heartbeat to server
                    this.heartbeat();
                    // 利用bind，解决this指针丢失的问题
                    // https://blog.csdn.net/Victor2code/article/details/107804354
                    this.heartbeatInterval = setInterval(this.heartbeat.bind(this), 30 * 1000);
                    break;
                case OpHeartbeatReply:
                    console.log('receive heartbeat reply');
                    appendToDialog('server: reply heartbeat');
                    break;
                case OpBatchMsg:
                    // batch message
                    for (let offset = rawHeaderLen; offset < data.byteLength; offset += packetLen) {
                        let packetLen = dataView.getInt32(offset);
                        let headerLen = dataView.getInt16(offset + headerOffset);
                        let ver = dataView.getInt16(offset + verOffset);
                        let op = dataView.getInt32(offset + opOffset);
                        let seq = dataView.getInt32(offset + seqOffset);
                        let msgBody = this.textDecoder.decode(data.slice(offset + headerLen, offset + packetLen));
                        appendToDialog("receive: ver=" + ver + " op=" + op + " seq=" + seq + " message=" + msgBody);
                    }
                    break;
                default:
                    let msgBody = this.textDecoder.decode(data.slice(headerLen, packetLen));
                    appendToDialog("receive: ver=" + ver + " op=" + op + " seq=" + seq + " message=" + msgBody);
                    break
                    console.log(event);
            }
        }
        /**
         * 当WebSocket对象的readyState状态变为CLOSED时会触发该事件。
         * 该事件表明这个连接已经已经关闭。
         * event属于CloseEvent对象，https://developer.mozilla.org/en-US/docs/Web/API/CloseEvent
         */
        this.WsClient.onclose = (event) => {
            appendToDialog("连接已关闭...");
            appendToDialog(event.reason);
            console.log(event);

        };
        /**
         * 当WebSocket发生错误时的回调。
         * event属于Event对象，https://developer.mozilla.org/en-US/docs/Web/API/Event
         */
        this.WsClient.onerror = (event) => {
            appendToDialog("连接出错...");
            console.log(event);
        }
    }

    // 关闭连接
    closeLink() {
        this.WsClient.close();
        this.WsClient = null;
    }

    // 发送消息
    sendMessage() {
        let msg = {
            'type': 'php',
            'msg': $("#msg-txt").val()
        }
        this.WsClient.send(JSON.stringify(msg));
    }

    heartbeat() {
        let headerBuf = new ArrayBuffer(rawHeaderLen);
        let headerView = new DataView(headerBuf, 0);
        headerView.setInt32(packetOffset, rawHeaderLen);
        headerView.setInt16(headerOffset, rawHeaderLen);
        headerView.setInt16(verOffset, 1);
        headerView.setInt32(opOffset, 2);
        headerView.setInt32(seqOffset, 1);
        // ws.WsClient.send(headerBuf);// lose this
        this.WsClient.send(headerBuf);// lose this
        console.log(this)
        console.log("send heartbeat to server");
        appendToDialog("client: send heartbeat");
    }

    // 授权
    auth() {
        let authInfo = '{"user_id":1001, "user_key":"x4u5mmq6gh2md5dl","room_id":"live://9999", "platform":"web", "token":"abcabcabcabc"}'
        let headerBuf = new ArrayBuffer(rawHeaderLen);
        let headerView = new DataView(headerBuf, 0);
        let bodyBuf = this.textEncoder.encode(authInfo);
        headerView.setInt32(packetOffset, rawHeaderLen + bodyBuf.byteLength);
        headerView.setInt16(headerOffset, rawHeaderLen);
        headerView.setInt16(verOffset, protoVersion);
        headerView.setInt32(opOffset, OpAuth);
        headerView.setInt32(seqOffset, 1);
        this.WsClient.send(this.mergeArrayBuffer(headerBuf, bodyBuf));
        appendToDialog("client: send auth" + authInfo + ".");
    }

    mergeArrayBuffer(ab1, ab2) {
        let u81 = new Uint8Array(ab1),
            u82 = new Uint8Array(ab2),
            res = new Uint8Array(ab1.byteLength + ab2.byteLength);
        res.set(u81, 0);
        res.set(u82, ab1.byteLength);
        return res.buffer;
    }
}

let ws = new (WebsocketOp)
