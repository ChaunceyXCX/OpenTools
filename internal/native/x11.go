//go:build linux

package native

/*
#cgo linux pkg-config: gtk+-3.0
#cgo LDFLAGS: -lX11

#include <gtk/gtk.h>
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <gdk/gdkx.h>

static void apply_ewmh(GtkWindow *win) {
	GdkWindow *gdk_window = gtk_widget_get_window(GTK_WIDGET(win));
	if (gdk_window == NULL) return;
	GdkDisplay *display = gdk_window_get_display(gdk_window);
	Display *xdisplay = gdk_x11_display_get_xdisplay(display);
	Window xwindow = gdk_x11_window_get_xid(gdk_window);

	Atom type_atom = XInternAtom(xdisplay, "_NET_WM_WINDOW_TYPE", False);
	Atom popup = XInternAtom(xdisplay, "_NET_WM_WINDOW_TYPE_POPUP_MENU", False);
	XChangeProperty(xdisplay, xwindow, type_atom, XA_ATOM, 32,
		PropModeReplace, (unsigned char *)&popup, 1);

	Atom state_atom = XInternAtom(xdisplay, "_NET_WM_STATE", False);
	Atom skip_tb = XInternAtom(xdisplay, "_NET_WM_STATE_SKIP_TASKBAR", False);
	Atom skip_pg = XInternAtom(xdisplay, "_NET_WM_STATE_SKIP_PAGER", False);
	Atom above = XInternAtom(xdisplay, "_NET_WM_STATE_ABOVE", False);
	Atom states[] = {skip_tb, skip_pg, above};
	XChangeProperty(xdisplay, xwindow, state_atom, XA_ATOM, 32,
		PropModeReplace, (unsigned char *)&states, 3);
	XFlush(xdisplay);
}

static void on_realize(GtkWidget *widget, gpointer data) {
	apply_ewmh(GTK_WINDOW(widget));
}

void set_floating(GtkWindow *win) {
	gtk_widget_realize(GTK_WIDGET(win));
	apply_ewmh(win);
	g_signal_connect_after(win, "realize", G_CALLBACK(on_realize), NULL);
}
*/
import "C"
import "unsafe"

func SetFloatingWindow(nativeWindow unsafe.Pointer) {
	if nativeWindow == nil {
		return
	}
	C.set_floating((*C.GtkWindow)(nativeWindow))
}
