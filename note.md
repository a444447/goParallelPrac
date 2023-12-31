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

