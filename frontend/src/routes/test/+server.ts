export function GET() {
    const now = new Date()
    let message = `data: This is a message at time ${now.toLocaleTimeString("de")}\n\n`
    const sleep = (milliseconds: number) => {
        return new Promise(resolve => setTimeout(resolve, milliseconds))
    }
    while (true) {
        sleep(100)
        if (now.getSeconds() % 2 == 0) {
            message = 'event: userconnect\ndata: {"username": "asd", "time": "02:33:48"}\n\n'
        } else {
            message = 'event: usermessage\ndata: {"username": "tim", "time": "asdasd"}\n\n'
        }
        return new Response(message, {
            headers: {
                'X-Accel-Buffering': 'no',
                'Content-Type': 'text/event-stream',
                'Cache-Control': 'no-cache'
            }
        });
    }
}