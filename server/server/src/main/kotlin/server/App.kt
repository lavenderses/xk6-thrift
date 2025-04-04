package server

import com.linecorp.armeria.common.HttpResponse
import com.linecorp.armeria.common.thrift.ThriftSerializationFormats
import com.linecorp.armeria.server.Server
import com.linecorp.armeria.server.docs.DocService
import com.linecorp.armeria.server.logging.AccessLogWriter
import com.linecorp.armeria.server.thrift.THttpService
import idl.Feature
import idl.Message
import idl.TestService
import org.apache.thrift.async.AsyncMethodCallback
import org.slf4j.LoggerFactory

class App {
    companion object {
        val log = LoggerFactory.getLogger(App::class.java)
    }
}

fun main() {
    Server.builder()
        .http(8080)
        .service("/health", { _, _ -> HttpResponse.of("OK")})
        .service("/test") { _, _ -> HttpResponse.of("Hello, world!") }
        .serviceUnder("/docs", DocService())
        .accessLogWriter(AccessLogWriter.common(), true)
        .service(
            "/thrift",
            THttpService.ofFormats(
                /* implementation = */ Service(),
                /* defaultSerializationFormat = */ ThriftSerializationFormats.BINARY,
                /* otherSupportedSerializationFormats = */ listOf(
                    ThriftSerializationFormats.JSON,
                    ThriftSerializationFormats.TEXT,
                    ThriftSerializationFormats.COMPACT,
                ),
            )
        )
        .build()
        .apply {
            closeOnJvmShutdown().thenRun {
                App.log.info("Server has been stopped.")
            }

            start().join()
        }
}

class Service: TestService.AsyncIface {
    companion object {
        private val log = LoggerFactory.getLogger(App::class.java)
    }

    override fun simpleCall(id: String, resultHandler: AsyncMethodCallback<String>) {
        // Make failure if specified
        if (id == "FAILURE") {
            val msg = "Make failure: $id"
            log.error(msg)
            resultHandler.onError(RuntimeException(msg))
        } else {
            val msg = "Success: $id"
            log.info(msg)
            resultHandler.onComplete(msg)
        }
    }

    override fun boolCall(tf: Boolean, resultHandler: AsyncMethodCallback<Boolean>) {
        log.info("Accept: $tf")
        resultHandler.onComplete(tf)
    }

    override fun messageCall(message: Message, resultHandler: AsyncMethodCallback<Message>) {
        log.info("Accept: $message")
        val newMessage = Message("content: ${message.content}", message.tags, message.nested)
        log.info("Respond: $newMessage")
        resultHandler.onComplete(newMessage)
    }

    override fun mapCall(
        maps: MutableMap<String, Boolean>,
        resultHandler: AsyncMethodCallback<MutableMap<String, Boolean>>,
    ) {
        log.info("Accept: $maps")
        val newMaps = maps.map { "NEW: ${it.key}" to it.value }.toMap().toMutableMap()
        log.info("Respond: $newMaps")
        resultHandler.onComplete(newMaps)
    }

    override fun stringCall(strs: MutableList<String>, resultHandler: AsyncMethodCallback<MutableList<String>>) {
        log.info("Accept: $strs")
        val response = strs.map { "$it:$it" }
        log.info("Respond: $response")
        resultHandler.onComplete(response.toMutableList())
    }

    override fun stringsCall(strs: MutableList<Message>, resultHandler: AsyncMethodCallback<MutableList<Message>>) {
        log.info("Accept: $strs")
        val response = strs.map { Message("content: ${it.content}", it.tags, it.nested) }
        log.info("Respond: $response")
        resultHandler.onComplete(response.toMutableList())
    }

    override fun enumCall(feature: Feature, resultHandler: AsyncMethodCallback<MutableList<Feature>>) {
        log.info("Accept: $feature")
        val response = listOf(Feature.ONE, Feature.TWO, Feature.THREE)
        log.info("Respond: $response")
        resultHandler.onComplete(response.toMutableList())
    }
}
