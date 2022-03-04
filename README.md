## 微信公众号：SuperGopher

go、云原生技术、项目问题、单纯支持 ...... 来者不拒
<div align="center" style="width: 50%">

![微信公众号.png](https://user-images.githubusercontent.com/55381228/155444889-eacc0104-cd85-45c9-b7b7-9036e0c2334c.jpg)
</div>

## 起源

每个客户端开发者都会想独立开发一款自己的 APP 。 但 iOS 不像 Android 那样可以自由分发应用 （Android 只要把 apk 甩出去就行了）

对于 iOS 开发者来说，苹果开发者账号几乎人手必备（这里只讨论个人账号）， 而苹果公司允许我们添加 100 台设备（udid）绑定到账号上，这 100 台设备可以自由安装由账号签名且使用 Ad Hoc 方式打包出的 `.ipa`

本项目就是利用这个规则，来简化 iOS APP 的分发流程。想一想，当你开发一款 APP 的过程中，想要给身边的小伙伴体验一下，只需要使用本项目生成一个二维码链接，扫一扫，即可全程自动绑定设备并签名安装 APP 。

## 声明

本项目的核心功能调用 [zsign](https://github.com/zhlynn/zsign)
和 [App Store Connect API](https://developer.apple.com/documentation/appstoreconnectapi) 实现，如有侵权，请与我联系，会及时处理。

本项目仅作为给开发者分发合法合规的 APP 使用，严禁使用本项目进行任何盈利、损害官方利益、分发任何违法违规的 APP 等行为。

## 部署项目

待更

## 效果预览

待更

## 喝杯奶茶

如有帮助，可以打赏支持，一分也是爱！

|  ![微信打赏](https://user-images.githubusercontent.com/55381228/155450359-0ce92911-fd3f-4d6b-9878-e40a17b34652.jpg)   | ![支付宝打赏](https://user-images.githubusercontent.com/55381228/155450383-509d0475-5497-4983-8583-137946b4d78e.jpg)  |
|  ----  | ----  |
| 微信  | 支付宝 |

## JetBrains 开源证书支持

本项目使用 GoLand 开发，感谢 JetBrains 提供的免费授权

<a href="https://www.jetbrains.com/?from=togettoyou" target="_blank"><img src="https://user-images.githubusercontent.com/55381228/127271051-14879011-41dd-4d1b-88a2-1591925b51de.png" width="250" align="middle"/></a>
