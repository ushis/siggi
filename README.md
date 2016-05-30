# Siggi

WebSocket server suitable for WebRTC signaling.

It is used in production on https://gibbler.honkgong.info/

## Build

```
go get github.com/ushis/siggi
```

## Run

```
siggi -listen=:9000
```

## API

Siggi relays JSON encoded messages between peers in the same room.

To join a room open a WebSocket and specifiy the room.

```js
var socket = new WebSocket('ws://localhost:9000/socket?room=fancy-room-name');
```

To broadcast a message to all peers connected to the room, send a message
without an recipient.

```js
socket.send(JSON.stringify({
  type: 'hello',
  data: 'arbitrary data'
}));
```

To send a message to specific peer, add an recipient.

```js
socket.send(JSON.stringify({
  type: 'hello',
  to: 'recipient-id',
  data: 'arbitrary data'
}));
```

The recipient will receive the following message.

```js
{
  type: 'hello',
  from: 'sender-id',
  to:   'recipient-id',
  room: 'fancy-room-name',
  data: 'arbitrary data'
}
```

### Message Fields

| Name | Contents         | Set By |
|------|------------------|--------|
| type | arbitrary string | sender |
| from | sender id        | server |
| to   | recipient id     | sender |
| room | room id          | server |
| data | arbitrary object | sender |


## License (GPL v2)

```
Copyright (C) 2015 The siggi authors

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301,
USA.
```
