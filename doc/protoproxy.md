## Basic Type

**ADDR** represent a raw IP address

```
{
  INT(16) : Port number
  BYTE : Length of address (For IPv4 is 4 bytes, for IPv6 is 16 bytes),
  []BYTE : IP Address,
}
```

**ANY** is matched any defined types.

**BYTE** represent a byte.

**INT(LEN)** is a unsigned int with. *LEN* is it length in count of bits.

**PARAM** is a parameter

```
{
  BYTE : Parameter Type,
  []ANY : Parameter Data
}
```

**STRING** is a text

```
{
  INT(16) : LEN,
  [LEN]BYTE : text in UTF-8,
}
```

**UUID** represent a 16 bytes UUID

### Frame Header

```
{
  BYTE : Frame Type,
  INT(24) : Frame Size,
}
```

### Commands

#### Alive Test (0x00)

Body: None

#### General Reply (0x01)

Body:

```
{
  INT(16) : result code
}
```

Zero for succeed, otherwise for failed.

#### Create Session (0x02)

Create and switch to new session.

Body:

```
{
  BYTE : Session Type,
  UUID : Session UUID,
  BYTE : Count of Parameters,
  []PARAM : Parameters,
}
```

Allowed Session Type:

- Stream Connect (0x00)
- Stream Bind (0x01)
- Reverse Stream Accept (0x02)
- Datagram Bind(0x03)

Allowed Parameter Type:

- Target Address (0x00) -> STRING target Address of session
- Idle Timeout (0x01) -> INT(32) Duration in seconds
- Alive Timeout (0x02) -> INT(32) Duration in seconds

#### Select Session (0x03)

Switch to an existed session

Body:

```
{
  UUID : Session UUID,
}
```

#### Close Session (0x04)

Close current selected session and wait for close connection.

Body: None

#### Event Notification (0x05)

Body:

```
{
  BYTE : Event ID,
  BYTE : Count of Parameters,
  []PARAM : Parameters,
}
```

See *Events* Chapter

#### Payload Stream (0x60)

Body:

```
{
  []BYTE : payload data,
}
```

#### Payload Datagram (0x61)

Body:

```
{
  ADDR : Remote Address,
  []BYTE : payload data,
}
```

### Events

#### Text Notification(0x00)

Text message between server

Body:

```
{
  STRING : Message,
}
```

#### Session Closed(0x01)

Report sesson have been closed.

Body:

```
{
  UUID : Session UUID,
  BYTE : Error code,
}
```

#### Reverse Stream Request(0x02)

Request to accept inbound reverse stream.

Body:

```
{
  UUID : Session UUID,
}
```

#### Datagram Request(0x03)

Request to receive inbound datagram while session not been selected.

Body:

```
{
  UUID : Session UUID,
}
```

#### Application Notification (0x80)

Custom notification.

Body:

```
{
  INT(32) : Custom notification code,
  STRING : Message,
}
```
