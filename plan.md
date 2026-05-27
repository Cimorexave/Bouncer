## components:
1. processor: this will create an imaginary process where it sends requests to the input channel. this part can later be expanded with feaetures for testing by enabling configurations for it to setup the initial load, load intervals, and etc.
2. channels: channels are created on the main app where requests are sent into them. this includes an input channel where all the routines use to read from and an eventual output channel where the responses go into them and wait to be sent to the next layer
3. routines: routines are the main component of this structure. the main routines will process each coming request and filter them based on their ip address and if that request gets flagged because it was too frequent in the input channel they are pushed into a wait queue. if not they are passed into the output where it passes to the next layer.
some other routines work on the wait queue. they will constantly look for old requests and force push them into main queue as well as when the signal is received when the workload from main queue is not so busy they will start slowly trickle into the main queue.
4. queues: includes of main and wait queue. they are both considered as a pipe data structure when the first one out is last one in (FIFO). each main routine is tied to a go routine. when they scale, these scale as well. (point of discussion?)

## Update
- queues are basically of the type of channels. under heavy workload in this priority the architecture handles:
increase channels buffer size -> scale workers -> scale channels
- routines have different type. there is an original routine where it handles the logic and its work is bounded by cpu performance. this worker handles data between main queue channel and output channel. where load handling routines of R1, R2, and R3 are working between the main and wait queue channels.