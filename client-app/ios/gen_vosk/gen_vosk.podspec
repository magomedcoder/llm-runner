Pod::Spec.new do |s|
  s.name             = 'gen_vosk'
  s.version          = '0.1.0'
  s.summary          = 'Vosk C API (XCFramework) для Flutter FFI'
  s.description      = <<-DESC
    и https://github.com/alphacep/vosk-api/tree/master/ios
  DESC
  s.homepage         = 'https://github.com/alphacep/vosk-api'
  s.license          = { :type => 'Apache-2.0', :file => 'LICENSE' }
  s.author           = { 'gen' => 'gen' }
  s.source           = { :path => '.' }
  s.platform         = :ios, '13.0'
  s.requires_arc     = false

  fw = File.join(__dir__, 'Vosk.xcframework')
  unless File.directory?(fw)
    raise "#{s.name}: отсутствует #{fw}. Запустите pod install с GEN_IOS_VOSK=1 (ensure_vosk_ios.sh)."
  end

  s.static_framework = true
  s.vendored_frameworks = 'Vosk.xcframework'
end
