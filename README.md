Queue
=====

When using the go language, we often use bytes.Buffer for byte stream reading and writing. Can I have a queue like bytes.Buffer?It can read and write to stream of any type stream. That's why I wrote this queue. In fact most of its code is copied from bytes.Buffer. 

The queue is not thread-safe, thread-safe is handled by the caller himself
