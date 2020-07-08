// This is simple push receiver service worker.
// Example not used because it doesn't allow customization (setBackgroundEventHandler)
self.addEventListener("push", event => { event.waitUntil(onPush(event)) });

async function onPush(event) {
    const push = event.data.json();
    console.log("push received", push)

    const { notification = {}, data = {} } = {...push};

    await self.registration.showNotification(notification.title, {
        body: notification.body,
    })

    if (data.id) {
        await fetch(`${self.location.origin}/api/v1/notifications/${data.id}/confirm`, { method: "POST" })
    }
}
