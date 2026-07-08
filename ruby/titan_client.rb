require 'socket'

class TitanClient
  CMD_PUT    = 1
  CMD_GET    = 2
  CMD_DELETE = 3
  CMD_3SUM   = 4

  def initialize(host = '127.0.0.1', port = 8001)
    @host, @port = host, port
    @socket = nil
    @message_id = 0
  end

  def connect
    @socket = TCPSocket.new(@host, @port)
  rescue Errno::ECONNREFUSED
    raise "Não foi possível conectar ao TitanDB em #{@host}:#{@port}"
  end

  def put(key, value)
    payload = [key.bytesize, key, value.bytesize, value].pack("S>a*S>a*")
    send_frame(CMD_PUT, payload)
    read_frame
  end

  def get(key)
    send_frame(CMD_GET, key)
    read_frame
  end

  def delete(key)
    send_frame(CMD_DELETE, key)
    read_frame
  end

  def three_sum
    send_frame(CMD_3SUM, "")
    read_frame
  end

  def close
    @socket&.close
    @socket = nil
  end

  private

  def send_frame(cmd, payload)
    connect if @socket.nil? || @socket.closed?
    @message_id += 1
    header = [payload.bytesize, @message_id, cmd].pack("L>L>C")
    @socket.write(header + payload)
  rescue Errno::EPIPE, Errno::ECONNRESET
    close
    raise "Conexão perdida com o TitanDB."
  end

  def read_frame
    header_raw = @socket.read(9)
    return nil if header_raw.nil? || header_raw.bytesize < 9
    payload_size, msg_id, status = header_raw.unpack("L>L>C")
    payload_data = payload_size > 0 ? @socket.read(payload_size) : ""
    { message_id: msg_id, status: status, data: payload_data }
  end
end