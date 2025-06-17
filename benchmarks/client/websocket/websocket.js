/**
 * 输出信息到 textarea
 * @param {string} content
 */
function appendToDialog(content) {
    const logEl = document.getElementById('log');
    if (!logEl) return;

    const timestamp = date();
    logEl.value += `${timestamp} | ${content}\r\n`;
    logEl.scrollTop = logEl.scrollHeight;
}


/**
 * 返回当前时间字符串
 * @returns {string}
 */
function date() {
    const today = new Date();
    const date = `${today.getFullYear()}-${today.getMonth() + 1}-${today.getDate()}`;
    const time = `${padZero(today.getHours())}:${padZero(today.getMinutes())}:${padZero(today.getSeconds())}`;
    return `${date} ${time}`;
}

/**
 * 补零函数
 * @param {number} num
 * @returns {string}
 */
function padZero(num) {
    return num.toString().padStart(2, '0');
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

// op code
const OpHeartbeat = 2
const OpHeartbeatReply = 3
const OpSendMsg = 4
const OpAuth = 7
const OpAuthReply = 8
const OpBatchMsg = 9

class WebsocketOp {
    wsClient;
    seqNum;

    constructor() {
        this.wsClient = null;
        this.seqNum = 1;
        this.textEncoder = new TextEncoder();
        this.textDecoder = new TextDecoder();
    }

    reset() {
        this.seqNum = 1
        clearInterval(this.heartbeatInterval)
    }

    // 连接
    connect() {
        const url = document.getElementById("ws-url").value;
        if (this.wsClient && this.wsClient.readyState === WebSocket.OPEN) {
            appendToDialog("已连接，无需重复连接。");
            return;
        }
        this.wsClient = new WebSocket(url);
        this.wsClient.binaryType = 'arraybuffer';
        this.reset()
        /**
         * 当WebSocket对象的readyState状态变为OPEN时会触发该事件。
         * 该事件表明websocket连接成功并开始发送数据。
         * event属于Event对象，https://developer.mozilla.org/en-US/docs/Web/API/Event
         */
        this.wsClient.onopen = (event) => {
            appendToDialog("连接成功...");
            console.log(event);
            this.sendAuth()
        };
        /**
         * 当有消息到达客户端的时候该事件会触发
         * event属于MessageEvent对象，https://developer.mozilla.org/en-US/docs/Web/API/MessageEvent
         */
        this.wsClient.onmessage = (event) => {
            console.log(event)
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
                    this.sendHeartbeat();

                    // 利用bind，解决this指针丢失的问题
                    // https://blog.csdn.net/Victor2code/article/details/107804354
                    this.heartbeatInterval = setInterval(this.sendHeartbeat.bind(this), 30 * 1000);
                    break;
                case OpHeartbeatReply:
                    console.log('receive heartbeat reply');
                    appendToDialog('server: reply heartbeat');
                    break;
                case OpBatchMsg:
                    // batch message
                    // 因为在switch之前已经解过一次包，所以offset的值从rawHeaderLen开始
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
                    console.log(event);
                    break
            }
        }
        /**
         * 当WebSocket对象的readyState状态变为CLOSED时会触发该事件。
         * 该事件表明这个连接已经已经关闭。
         * event属于CloseEvent对象，https://developer.mozilla.org/en-US/docs/Web/API/CloseEvent
         */
        this.wsClient.onclose = (event) => {
            this.reset()
            appendToDialog("连接已关闭...");
            appendToDialog("event=" + event);
            console.log(event);

        };
        /**
         * 当WebSocket发生错误时的回调。
         * event属于Event对象，https://developer.mozilla.org/en-US/docs/Web/API/Event
         */
        this.wsClient.onerror = (event) => {
            this.reset()
            appendToDialog("连接时遇到错误...");
            appendToDialog("event=" + event);
            console.log(event);
        }
    }

    // 关闭连接
    disconnect() {
        this.wsClient.close();
        this.wsClient = null;
    }

    // 发送授权请求
    sendAuth() {
        const authInfo = JSON.stringify({
            user_info: {
                tcp_session_id: {
                    user_id: 1001,
                    user_key: 'x4u5mmq6gh2md5dl',
                },
                room_id: 'live://9999',
                platform: 4,
            },
            token: 'abcabcabcabc',
        });
        // send frame
        this.sendFrame(OpAuth, authInfo)
    }

    // 发送消息
    sendMessage() {
        const msgInput = document.getElementById('msg-txt');
        if (!msgInput) return;
        // send frame
        this.sendFrame(OpSendMsg, msgInput.value);
    }

    // 发送心跳
    sendHeartbeat() {
        // send frame
        this.sendFrame(OpHeartbeat, '')
        console.log("send heartbeat to server");
        appendToDialog("client: send heartbeat");
    }

    // 客户端发送消息
    sendFrame(op, body) {
        let headerBuf = new ArrayBuffer(rawHeaderLen);
        let headerView = new DataView(headerBuf, 0);
        let bodyBuf = this.textEncoder.encode(body);
        // length
        let totalLen = rawHeaderLen + bodyBuf.byteLength
        let headerLen = rawHeaderLen
        // pack
        headerView.setInt32(packetOffset, totalLen);
        headerView.setInt16(headerOffset, headerLen);
        headerView.setInt16(verOffset, protoVersion);
        headerView.setInt32(opOffset, op);
        headerView.setInt32(seqOffset, this.seqNum++);
        // send
        this.wsClient.send(this.mergeArrayBuffer(headerBuf, bodyBuf));
        appendToDialog("client: send frame: " + body + ".");
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