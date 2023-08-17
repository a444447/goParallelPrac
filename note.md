[toc]

## understanding parallel computing

一个重要的概念叫做*amdalslo's law* (阿姆达尔定律)

> amdalslo定律是主要是来观测提升一个系统一部分的效率，对整体程序性能的影响。
>
> 一般的定义是这样的，假设在一个系统里，执行一个程序所需要的总时间为$T_{total}$，其中一段程序所用时间占比为$\alpha$，我们将这段程序性能提升k倍，即这一部分原来需要的时间为 ![\alpha T_{old}](https://www.zhihu.com/equation?tex=%5Calpha+T_%7Bold%7D&consumer=ZHI_MENG) ，现在需要的时间变为 ![(\alpha T_{old})/k](https://www.zhihu.com/equation?tex=%28%5Calpha+T_%7Bold%7D%29%2Fk&consumer=ZHI_MENG) 。则整个系统执行此程序需要的时间变为：
>
> ![T_{new}=(1-\alpha)T_{old} + (\alpha T_{old})/k =T_{old}[(1-\alpha) + \alpha /k]](https://www.zhihu.com/equation?tex=T_%7Bnew%7D%3D%281-%5Calpha%29T_%7Bold%7D+%2B+%28%5Calpha+T_%7Bold%7D%29%2Fk+%3DT_%7Bold%7D%5B%281-%5Calpha%29+%2B+%5Calpha+%2Fk%5D&consumer=ZHI_MENG)
>
> 故可得，系统性能提速的倍数为: ![S=\frac{1}{(1-\alpha)+\alpha/k}](https://www.zhihu.com/equation?tex=S%3D%5Cfrac%7B1%7D%7B%281-%5Calpha%29%2B%5Calpha%2Fk%7D&consumer=ZHI_MENG)
>
> 现在对于并行和串行的程序，可以这样理解，$\alpha$表示的是程序里必须串行执行的比例，$t_1,t_2$分别表示优化前后的耗时，n为总共的处理器，那么有,$t_2=t_1\times \alpha + (1-\alpha)/n \times t_1$
> 加速的比例=t1/t2
>
> 如果串行比例占2/3，则无论处理器再多，最大加速比也只能达到1.5。理论上并行越多，加速比越大。

除了上面提到的定律外，还有一个定律叫做古斯塔夫森定律

> 古斯塔夫森定律也是在表明处理器个数、并行比例和加速比之间的关系。执行时间： 串行时间a + 并行时间b
>
> 优化后时间： a + nb、
>
> 加速比： (a + nb) / (a + b)
>
> f串行比例 ： a / (a + b)
>
> 加速比:$\frac{a+nb}{a+b}=\frac{a}{a+b} + n\frac{a+b-a}{a+b}=f + n(1-\frac{a}{a+b})=f+n-nf=n-f(n-1)$
>
> 如果串行比例很小，那个加速比就是处理器的个数。
>
> 两个定律最低点、最高点都是一致的结论：
>
> 无可并行的程序，加速比就是1.
> 全部是并行程序，加速比就是n。

即使是没有多核处理器的系统，单核也可以基于轮转等机制来达到并行的效果。



## process, threads, greenthreads

首先来看进程，以记事本为例，我们可以运行多个记事本程序，每个记事本互相都是独立的，它们有自己的内存空间，如果其中一个崩溃了，并不会影响另外的，这就是一种隔离(isolation)状态。

但是，进程是重量级的，因为它需要有自己的内存空间，因此需要消耗很大的资源。同样，创建一个进程也很消耗时间，因为需要作分配空间等一系列操作。

在golang中，更集中关注在threads和greenthreads

在创建线程的时候，我们就不需要单独分配内存空间，创立的线程们共享同一个内存空间，因此它是轻量级的。

> 以一张画布类比，假如我们想画一幅画，我们一共有三个人，进程就好比每个人都在一张独立的画布上作画，然后最后粘连组合为一张完整的。而线程就相当于给三支笔给三个人，让它们共同在一张画布画它们各自的部分。这节省了很多多余步骤，但是也需要了更多的沟通。

对于threads和greenthreads，一种更常见的称呼是*kernel threads以及user level threads*。讨论为什么引入用户级线程，我们首先要知道内核级线程的问题。思考这样一种场景，我们有一系列的进程排在等待队列中，让处理器处理，每当遇到IO型的任务或者是到达了一个分片的时间，就需要将当前任务移动到队列尾，操作系统选出下一个需要执行的任务。



操作系统切换作业，需要保存原来的信息，然后更新新的内存空间与程序计数器，这个过程称为 **上下文切换**。用户级线程致力于减少切换上下文所耗费的时间。使用用户级线程后，虽然对于操作系统来说，它只能看到内核级线程，但是对于程序来说，它是根据用户级线程来决定是否切换上下文。

用户级线程的问题在于，外部的操作系统对用户级线程一无所知，上面以及提到，操作系统只能看到内核级线程。因此当一个用户级线程需要进行IO操作时，它会将整个进程的都先阻塞，尽管这个进程中的其他用户级线程不是IO型的。

在golang中，采取了一个混合式的方法。它为每一个CPU核都创建了一个内核级线程，每一个这些内核级线程都有一些greenthreads.

当某个greenthreads做了IO型操作，那么整个线程都会被阻塞，此时其他的greenthreads将会被分配到其他的我们初始创建的线程上。

## using goroutines for boids

BOID=BIRD + ROBORT

我们使用Go来模拟这个过程，我们会模拟生成很多个BOID，每个都由一个单独的goroutine控制与其他BOIDS之间的交互。每个BOID具有的属性包括了，position(x and y), velocity(速度的表示是用一个速度向量表示，也就是用某个boid下个位置与当前位置的(x1-x2,y1-y2)), ID(用来区分boid)

### 实现思路

#### memory sharing between threads

线程之间有两种方法通信:

+ 第一种，通过传递消息(message pass)
+ 第二种，通过共享内存, 也就是一个公共的空间来交流(思考操作系统的内容)

使用哪种方法取决于问题，它们都有各自的优缺点。

---



我们用之前实现的boids来了解内存共享。在boids中，每个boid都要知道它半径范围内所有boids的位置与速度。

![image-20230802143346140](https://obsdian-1304266993.cos.ap-chongqing.myqcloud.com/typora/image-20230802143346140.png)

找到view radius的方法，最简单的想法是对每个boid,遍历计算其他所有boids与他的距离，选出小于view radius的，但是这种方法的时间复杂度很大O(n*n)

另一种方法是，把整个空间看作二维数组，起始时空间为空白，每个位置都填充-1，接下来对于有boids的位置填充boids的编号。我们引入一个Box(长宽为2*view radius)，保证所有小于view radius的boids都在里面。以后每次都只需要根据Boid更新盒子的位置就可以了。

![image-20230802144748235](https://obsdian-1304266993.cos.ap-chongqing.myqcloud.com/typora/image-20230802144748235.png)



#### using lock

在创建多线程的时候，使用锁是必要的，可以参考在操作系统中提到的`生产者-消费者问题`,`读者问题`,`哲学家用餐问题`.



我们之间实现的boidsMap,存在的一个就是如果不加锁的话，所有的线程都在同时对公共内存(也就是我们的boidsMap)作操作，会出现以下的问题:

+ 对同一个位置的重复计算平均velocity.
+ 某个位置的boid的值(x,y)，可能x正确，但是y被其他线程操作已经改变了。

因此我们需要进行加锁。

![image-20230803003254557](https://obsdian-1304266993.cos.ap-chongqing.myqcloud.com/typora/image-20230803003254557.png)

#### reader and writer lock

在实际编写代码的时候，我们将锁加在了两个地方，第一个地方就是读取我们的公共内存，计算出相关数据的时候，另一个就是更新值得时候，这也是写入操作。

问题在于，写入操作比读取操作快很多，因此如果每次有goroutine在读取公共内存就锁住，那么就会造成很多goroutines堵在一起等待前面的协程读取完毕。

当出现这种 **写入的时间远远低于读取时间的时候，为了提高效率，我们就可以使用读写锁**。可以参考操作系统中的内容，写锁只能在读锁是空的时候才能进入，读锁允许复数的协程一起进入，写锁只能一个一个使用。



## waiting group

使用waiting group，主要是控制我们main threads的流程，保证是在所有的被创建的子threads全部done过后，才继续主要流程。当每次创建一个goroutine时，waiting group计数器就会增加1，然后当子threads调用了done()表示已经完成的时候，计数器会减少1。



### 以一个文件查找为例子

现在以一个文件查找的项目为例子进行学习waiting group。

假如我们要调用查找文件的函数`search('/', 'cat.jpg')`,表示我们要在根目录下找`cat.jpg`,我们会走一个这样的算法: 如果找到，将其放入`RECORD`里面，然后继续在子目录找；如果没找到，继续在子目录找。可以看到，这是一个递归的过程，每遇到一个子目录就会执行一次递归。现在我们把递归用创建新goroutine代替，每次遇到子目录就创建一个新的goroutine,然后waiting group计数器加1，直到没有了子目录，执行done,计数器减少1。

![](https://obsdian-1304266993.cos.ap-chongqing.myqcloud.com/typora/image-20230814190829970.png)

## channel

我们之前提到过两个threads交流的方式，包括了共享内存与send message。对于共享内存，我们很难知道到底两方进行操作的是谁（比如有多对生产者与消费者的情况，我们不知道它们之间具体谁与谁交互）；对于发送message，如果1个threads要向100个threads交流，那样会很低效率。

channel，就是通信方式的共享内存。

```go
//定义一个channel,需要定义它接受的类型
ch = make(chan string)
```



---

考虑一个这样的问题，我们收集到了一批天气报告的数据，它们是以特殊的形式表达的。我们的程序想要做的是这样一件事：

1. 读取数据、
2. 特殊数据的格式转换
3. 提取特征
4. 将特征与我们实际任务要求结合。

可以看到，这4个步骤其实是线性的，也就是每个步骤都依赖上一个步骤。在这样的情况下我们如何利用并行提高效率呢？即使我们为每个步骤都分配一个threads，它们也不能在没有上一步结果的时候运行。这里，我们需要做的就是用 **pipeline**。所谓的pipline,也就是流水线，就是当它某一步骤完成后，它不会等待整个流程走完才读取下一个，而是处理好自己的部分后马上读取下一个部分。

在这里，我们用channel实现这个流水线，也就是负责每一部分的threads，从channel中取出未处理好的工作，待到处理完成后，打好标记，放回channel，而负责下一部分的threads检测到它所谓要的preprocessed数据已经被放入了相应channel后，马上取出处理，此时上一步的threads继续从channel中拿出未完成工作处理。

![image-20230815000802679](https://obsdian-1304266993.cos.ap-chongqing.myqcloud.com/typora/image-20230815000802679.png)



### threads pool

我们上面提到的channel默认的缓冲区大小是1，也就是只能给一拿一再放回这样的模式。

`channel`也可以设置缓冲区大小大于1。假如依然是有生产者与消费者两个threads, `buffer_size=3`，那么生产者可以一直生产直到`buffer_size`被装满，或者至少有一个被消费者利用才能解除阻塞；消费者只有当`buffer`里面有内容的时候才能从buffer拿去，否则会阻塞。

现在利用这样的channel,我们就可以实现线程池(threads pool)。threads pool机制是这样的: *threads pool会在开始阶段就准备很多threads, 其中有一个main threads,其他的threads都是workers,负责处理工作。当用户向应用发起请求，请求对象是main threads,然后它再寻找还空闲的threads,将请求交给它处理并返回结果。* 我们一般是会维护一个线程的buffer, 从中找到是否还有空闲的线程。

![image-20230815114327007](https://obsdian-1304266993.cos.ap-chongqing.myqcloud.com/typora/image-20230815114327007.png)

---

现在以这样一个项目为例子，计算多边形的面积。我们设计一个应用，用户请求计算多边形面积，我们的应用会安排一个空闲的线程来执行这个工作并且返回结果。

计算公式如下，注意，假设$a_1$是起点，按照顺时针方向，$a_n$是最后一个点，那么最后$a_n$的下一个点就是$a_1$.

![image-20230816101926234](https://obsdian-1304266993.cos.ap-chongqing.myqcloud.com/typora/image-20230816101926234.png)

```go
//我们想要匹配类似"(4,10),(12,8)..."这样的文本，用正则表达式
//使用regex.MustCompile
//FindAllStringSubmatch(text, n=-1)进行匹配。
//func(re *Regexp) FindAllStringSubmatch(s string, n int) [][]string
//FindAllStringSubmatch 是 FindStringSubmatch 的 'All' 版本；它返回表达式的所有连续匹配的切片，如包注释中的'All' 说明所定义。返回值 nil 表示不匹配。
m := "(3,4),(5,6)"
fmt.Println(r.FindAllStringSubmatch(m, -1))

//结果:
//[[(3,4) 3 4] [(5,6) 5 6]]
```

## condition variable

condition variable可以判断一个条件是否满足，不满足的话会堵塞在那里，直到条件满足。

```go
type Cond struct {
        // L is held while observing or changing the condition
        L Locker
        // contains filtered or unexported fields
}
type Cond
    func NewCond(l Locker) *Cond
    func (c *Cond) Broadcast()
    func (c *Cond) Signal()
    func (c *Cond) Wait()
type Locker

type Locker interface {
        Lock()
        Unlock()
}
A Locker represents an object that can be locked and unlocked.
```

+ `NewCond()`函数输入参数是一个`Locker`接口类型，即实现了锁功能的变量。

+ `Broadcast()`函数通知所有等待在`condition variable`的`goroutine`，
+ `Signal()`函数只会通知其中的一个`goroutine`。
+ `Wait()`会让`goroutine`阻塞在`condition variable`，等待条件成立。通常的做法是：

```go
c.L.Lock()
for !condition() {
    c.Wait()
}
... make use of condition ...
c.L.Unlock()
```

进入`Wait()`函数会解锁，离开`Wait()`函数会重新加锁。由于在“解锁->收到通知->重新加锁”这段逻辑中间有可能另一个同样`wait()`的`goroutine`抢先一步改变了条件，导致当前`goroutine`的条件不再成立，所以这块要使用循环检测。

### 矩阵计算的例子

现在设计一个计算矩阵乘法的程序，使用并行的方法加快速度。我们使用一个独立的线程来从文件中读取矩阵，其余线程分别计算每行的结果。

我们选择的应该是读写锁，因为我们有多个线程需要读取内存中数据来进行计算。在最开始的时候，我们的main threads(我们假定它是负责从文件中读取矩阵的)需要知道其他threads是否做好了准备，因此它会初始化一个`wait group`，等到所有的threads都`done`后，主线程中的`wait group`解锁。但是现在的问题是，其余的`worker threads`声明自己`done`后，就被堵塞了，而读写锁是不允许写锁在有读锁处于堵塞状态进行写入的。

为了解决这点，我们使用了`condition variable`，声明自己`wait()`，从而解锁，让main threads能够读取文件。

![image-20230816144833175](https://obsdian-1304266993.cos.ap-chongqing.myqcloud.com/typora/image-20230816144833175.png)

## deadlock

死锁的原因和操作系统中出现的原因一样。主要是资源或者说锁的竞争。

解决方法之一就是，把资源或者锁编号，然后以某种层次结构进行需求，也就是要求threads申请资源时规定它们先看什么，再看什么。

## barrier

barrier的主要作用是允许我们将不同的threads同步在一起。可以这样打比方: *barrier就是将处于不同位置的threads强制让他们在同一起跑线等着，等都到了起跑线，才让它们在进行*。

![image-20230816184225044](https://obsdian-1304266993.cos.ap-chongqing.myqcloud.com/typora/image-20230816184225044.png)

在go中，我们需要自己实现barrier结构，我们定义barrier结构如下

```go
type Barrier struct {
	count int
	cond  *sync.Cond
	mutex *sync.Mutex
	total int
}
```

我们需要实现一个`newBarrier`的函数来获得一个Barrier,实现一个`wait()`函数来达到等待其余threads到达barrier的功能。

对于第二个`wait`函数，我们的思路是这样的：当我们newBarrier的时候会传入threads数量，count与total就会被赋予这个值。当有一个threads抵达barrier，那么count就会减少一，当count已经为0但是又需要减少一的时候，说明了所有threads已经到达了Barrier，我们将count = total。另外我们需要mutex来保护中间的对count的操作，并且当某一threads到达barrier的时候，它需要释放锁，因此就需要cond.wait()。

```go
func (b *Barrier) Wait() {
	b.mutex.Lock()
	b.count--
	if b.count == 0 {
		b.count = b.total
		b.cond.Broadcast()
	} else {
		b.cond.Wait()
	}
	b.mutex.Unlock()
}
```

---

现在，重新来编写并行计算matrix乘法，这次我们不需要读写锁与condition variable,我们直接使用barrier实现。具体思路是这样的:*如下图，假如我们参与计算的线程有4个，依然是一个负责读数据的main threads，和其他三个负责计算的threads。我们引入两个barriers,分别是workStart与workComplete,起始时三个计算threads默认抵达workStart,需要等待main threads读取数据后到达workStart；此后main threads到达workComplete然后等待其余三个threads计算完毕后到达workComplete。之后重复上面的步骤*



![image-20230816204853230](https://obsdian-1304266993.cos.ap-chongqing.myqcloud.com/typora/image-20230816204853230.png)

## Atomic Variable

我们之前使用的是互斥锁来避免资源的竞争，但是每次使用互斥锁的时候，lock与unlock操作都会让操作系统执行逻辑，完成后又要回到线程中，这会耗费额外的时间，在那种非常要求高性能的并行环境可能表现不好。

而原子操作是无法被中断的，因此使用原子操作对数据操作是不会产生之间的问题的。

原子操作的意思是说，这个操作在执行的过程中，其它协程不会看到执行一半的操作结果。在其它协程看来，原子操作要么执行完成了，要么还没开始，就像一个原子一样，不可分割。

原子操作主要是两类：修改和加载存储。修改很好理解，就是在原来值的基础上改动；加载存储就是读写。

atomic 提供了 AddXXX、CompareAndSwapXXX、SwapXXX、LoadXXX、StoreXXX 等方法。

由于 Go 暂时还不支持泛型，所以很多方法的实现都很啰嗦，比如 AddXXX 方法，针对 int32、int64、uint32 基础类型，每个类型都有相应的实现。等 Go 支持泛型之后，相信 atomic 的 API 就会清爽很多。

**需要注意的是，atomic 的操作对象是地址，所以传参的时候，需要传变量的地址，不能传变量的值。**

+ add:Add 方法很好理解，对 addr 指向的值加上 delta。如果将 delta 设置成负值，加法就变成了减法。
+ cas(compare and swap): 这个方法会比较当前 addr 地址对应值是不是等于 old，等于的话就更新成 new，并返回 true，不等于的话返回 false。
+ swap:如果不需要比较，直接交换的话，也可以用 Swap 方法。
+ load: Load 方法会取出 addr 地址中的值
+ store: Store 方法会将一个值存到指定的 addr 地址中去.



## spining locks

spining lock实现与mutex lock一样的功能。但是我们之前已经说过，mutex lock在`lock`与`unlock`的时候会将控制权给操作系统，如果遇到频繁操作这一步的话，效率会有影响。spining lock的实现就是，试图直接让CPU一直运行，而不交给操作系统，主要是利用我们上面提到的atomic variable实现功能。

其实就是操作系统中提到的旋转锁，用`compare and swap`实现。go中的Cas操作与java中类似，都是借用了CPU提供的原子性指令来实现。CAS操作修改共享变量时候不需要对共享变量加锁，而是通过类似乐观锁的方式进行检查，本质还是不断的占用CPU 资源换取加锁带来的开销（比如上下文切换开销）。

![image-20230817155516217](https://obsdian-1304266993.cos.ap-chongqing.myqcloud.com/typora/image-20230817155516217.png)

`cas(old, new)`只有当内存旧值与old相同时，才会更新为new。当多个线程尝试使用CAS同时更新一个变量的时候，只有一个能够更新成功。那就是当我们的内存值V和旧的预期值A相等的情况下，才能将内存值V修改成B！然后失败的线程不会挂起，而是被告知失败，可以继续尝试（自旋）或者什么都不做！

