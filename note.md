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

